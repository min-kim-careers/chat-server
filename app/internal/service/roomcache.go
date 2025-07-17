package service

import (
	"chat-server/internal/cache"
	"context"

	"github.com/redis/go-redis/v9"
)

func (s *RoomService) GetRoomChannel(ctx context.Context, roomID string) *redis.PubSub {
	return s.c.Client.Subscribe(ctx, cache.CacheStreamKey(roomID))
}
