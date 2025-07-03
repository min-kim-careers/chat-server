package helper

import (
	"chat-server/internal/constant"
	"time"
)

func ConvertTimestamp(timestamp string) float64 {
	t, err := time.Parse(constant.TIME_STANDARD, timestamp)
	if err != nil {
		return 0
	}

	return float64(t.Unix())
}
