package service

import "fmt"

func msgKey(roomID string) string {
	return fmt.Sprintf("chat:room:%s:msgs", roomID)
}

func stagingKey(roomID string) string {
	return fmt.Sprintf("chat:room:%s:staging", roomID)
}
