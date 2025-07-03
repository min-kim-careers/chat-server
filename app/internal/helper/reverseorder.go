package helper

import (
	"chat-server/internal/dto"
)

func ReverseOrder(msgs []*dto.Message) []*dto.Message {
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs
}
