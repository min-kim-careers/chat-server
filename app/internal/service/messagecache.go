package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/dto/messageout"
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type CacheChatMessageParams struct {
	ClientID string
	RoomID   string
	Content  string
}

func (s *MessageService) CacheChatMessage(ctx context.Context, key string, v map[string]any) error {
	return s.c.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		Values: v,
	}).Err()
}

func (s *MessageService) PersistChatMessage(ctx context.Context, v map[string]any) error {
	return s.c.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: cache.PersistStreamKey(),
		Values: v,
	}).Err()
}

func (s *MessageService) CacheAndPersistChatMessage(ctx context.Context, arg CacheChatMessageParams) (*cache.CacheMessage, error) {
	m, err := cache.NewCacheMessage(arg.RoomID, arg.ClientID, arg.Content)
	if err != nil {
		return nil, err
	}

	v := cache.ToCacheMessageValue(m)

	err = s.CacheChatMessage(ctx, cache.CacheStreamKey(arg.RoomID), v)
	if err != nil {
		log.Println("error caching:", err)
	} else {
		log.Println("Cached")
	}

	err = s.PersistChatMessage(ctx, v)
	if err != nil {
		log.Println("error persisting:", err)
	} else {
		log.Println("Persisted")
	}

	log.Println(m)
	return m, nil
}

type PublishMessageParams struct {
	Mode     string
	ClientID string
	RoomID   string
	Content  string
}

func (s *MessageService) PublishMessage(ctx context.Context, roomID string, m messageout.MessageOut) error {
	p, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = s.c.Client.Publish(ctx, cache.CacheStreamKey(roomID), p).Err()
	if err != nil {
		return err
	}

	log.Println("Published:", string(p))
	return nil
}

type GetCachedChatMessagesParams struct {
	ClientID    string
	RoomID      string
	Limit       int64
	LastCacheID string
}

func (s *MessageService) GetCachedChatMessages(
	ctx context.Context,
	arg GetCachedChatMessagesParams,

) ([]*messageout.MessageOutChat, *string, error) {
	msgs, err := s.c.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{cache.CacheStreamKey(arg.RoomID), arg.LastCacheID},
		Count:   arg.Limit,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, nil, err
	}

	if len(msgs) == 0 {
		return []*messageout.MessageOutChat{}, nil, nil
	}

	var lastCacheID string
	dtos := make([]*messageout.MessageOutChat, len(msgs))
	for i, m := range msgs[0].Messages {
		dto, err := cacheToMessageOutChat(m.Values, arg.ClientID)
		if err != nil {
			return nil, nil, err
		}
		dtos[i] = dto
		if i == len(msgs[0].Messages) {
			lastCacheID = m.ID
		}
	}

	return dtos, &lastCacheID, nil
}

func (s *MessageService) GetChatMessagePersistStream(ctx context.Context, workerID string, roomID string) ([]redis.XStream, error) {
	msgs, err := s.c.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    cache.PersistGroupKey(),
		Consumer: workerID,
		Streams:  []string{cache.CacheStreamKey(roomID), ">"},
		Count:    10,
		Block:    0,
	}).Result()
	if err != nil {
		return []redis.XStream{}, err
	}
	return msgs, nil
}
