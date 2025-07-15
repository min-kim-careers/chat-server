package service

import (
	"chat-server/internal/cache"
	"chat-server/internal/dto/messageout"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CacheChatMessageParams struct {
	ClientID string
	RoomID   string
	Content  string
}

func (s *MessageService) CacheChatMessage(ctx context.Context, arg CacheChatMessageParams) (*cache.MessageCache, error) {
	newID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	_roomID, err := uuid.Parse(arg.RoomID)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	c := &cache.MessageCache{
		ID:        newID,
		Mode:      "chat",
		RoomID:    _roomID,
		ClientID:  arg.ClientID,
		Content:   arg.Content,
		CreatedAt: t,
		Read:      false,
		Sent:      true,
	}
	p, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	err = s.c.Client.ZAdd(ctx, msgKey(arg.RoomID), redis.Z{
		Score:  float64(t.UnixNano()),
		Member: p,
	}).Err()
	if err != nil {
		return nil, err
	}

	log.Println("Cached:", c)
	return c, nil
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
	err = s.c.Client.Publish(ctx, msgKey(roomID), p).Err()
	if err != nil {
		return err
	}

	log.Println("Published:", string(p))
	return nil
}

type GetCachedChatMessagesParams struct {
	ClientID string
	RoomID   string
	Before   time.Time
	Limit    int64
}

func (s *MessageService) GetCachedMessages(
	ctx context.Context,
	arg GetCachedChatMessagesParams,

) ([]*messageout.MessageOutChat, error) {
	rows, err := s.c.Client.ZRevRangeByScore(ctx, msgKey(arg.RoomID), &redis.ZRangeBy{
		Max:   fmt.Sprintf("(%d", arg.Before.UnixNano()),
		Min:   "-inf",
		Count: arg.Limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []*messageout.MessageOutChat{}, nil
	}

	dtos := make([]*messageout.MessageOutChat, len(rows))
	for i, r := range rows {
		dto, err := cacheToMessageOut(r, arg.ClientID)
		if err != nil {
			return nil, err
		}
		dtos[i] = dto
	}

	return dtos, nil
}

func (s *MessageService) GetCachedChatMessagesSize(ctx context.Context, roomID string) int64 {
	card, err := s.c.Client.ZCard(ctx, msgKey(roomID)).Result()
	if err != nil {
		log.Println("error getting cache size:", err)
		return -1
	}
	return card
}
