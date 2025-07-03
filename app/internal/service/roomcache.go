package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func roomKey(roomID string) string {
	return "chat:room:" + roomID
}

func (s *RoomService) GetRoomChannel(ctx context.Context, roomID string) *redis.PubSub {
	return s.c.PubSub(ctx, roomKey(roomID))
}
