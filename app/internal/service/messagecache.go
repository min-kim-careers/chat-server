package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/constants"
	"chat-server/internal/dto/messageout"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
		log.Println("error:", err)
		return nil, err
	}

	v := cache.ToCacheMessageValue(m)

	err = s.CacheChatMessage(ctx, cache.CacheStreamKey(arg.RoomID), v)
	if err != nil {
		log.Println("error:", err)
	} else {
		log.Println("cached")
	}

	err = s.PersistChatMessage(ctx, v)
	if err != nil {
		log.Println("error:", err)
	} else {
		log.Println("persisted")
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
		log.Println("error:", err)
		return err
	}
	err = s.c.Client.Publish(ctx, cache.CacheStreamKey(roomID), p).Err()
	if err != nil {
		log.Println("error:", err)
		return err
	}

	log.Println("published:", string(p))
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
	msgs, err := s.c.Client.XRangeN(ctx, cache.CacheStreamKey(arg.RoomID), fmt.Sprintf("(%s", arg.LastCacheID), "+", constants.RESTORE_LIMIT).Result()
	if err != nil {
		log.Println("error:", err)
		return nil, nil, err
	}

	if len(msgs) == 0 {
		return []*messageout.MessageOutChat{}, nil, nil
	}

	var lastCacheID string
	dtos := make([]*messageout.MessageOutChat, len(msgs))
	for i, m := range msgs {
		dto, err := cacheToMessageOutChat(m.Values, arg.ClientID)
		if err != nil {
			log.Println("error:", err)
			return nil, nil, err
		}
		dtos[i] = dto
		lastCacheID = m.ID
	}

	log.Printf("fetched %d from cache", len(dtos))
	return dtos, &lastCacheID, nil
}

func (s *MessageService) GetCacheMessageStream(ctx context.Context, count int64, block time.Duration) ([]redis.XStream, error) {
	res, err := s.c.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    cache.PersistGroupKey(),
		Consumer: "persist_reader",
		Streams:  []string{cache.PersistStreamKey(), ">"},
		Count:    count,
		Block:    block,
	}).Result()
	if err != nil && err != redis.Nil {
		log.Println("error:", err)
		time.Sleep(time.Second)
		return nil, err
	}
	return res, nil
}

func (s *MessageService) MarkCacheBatchAsPersisted(ctx context.Context, msgIDs []string) error {
	return s.c.Client.XAck(ctx, cache.PersistStreamKey(), cache.PersistGroupKey(), msgIDs...).Err()
}

func (s *MessageService) GetPendingCacheMesages(ctx context.Context) (*redis.XPending, error) {
	return s.c.Client.XPending(ctx, cache.PersistStreamKey(), cache.PersistGroupKey()).Result()
}
