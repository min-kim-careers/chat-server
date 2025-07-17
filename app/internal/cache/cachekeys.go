package cache

import "fmt"

func CacheStreamKey(roomID string) string {
	return fmt.Sprintf("cache_stream:%s", roomID)
}

func PersistStreamKey() string {
	return "persist_stream"
}

func PersistGroupKey() string {
	return "persist_group"
}

func ReadKey(roomID string, clientID string) string {
	return fmt.Sprintf("chat_read:%s:%s", roomID, clientID)
}

func SentKey(roomID string, clientID string) string {
	return fmt.Sprintf("chat_sent:%s:%s", roomID, clientID)
}
