package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"chat-server/internal/models"
)

var ISOTimestampLayout = "2006-01-02T15:04:05.999Z"

func partition(arr []*models.Message, low, high int) ([]*models.Message, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		t1, _ := time.Parse(ISOTimestampLayout, arr[j].Timestamp)
		t2, _ := time.Parse(ISOTimestampLayout, pivot.Timestamp)
		if t1.Before(t2) {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []*models.Message, low, high int) []*models.Message {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func QuickSortStart(msgs []*models.Message) []*models.Message {
	return quickSort(msgs, 0, len(msgs)-1)
}

func ReverseOrder(msgs []*models.Message) []*models.Message {
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs
}

func ConvertTimestamp(timestamp string) float64 {
	t, err := time.Parse(TIMESTAMP_FORMAT, timestamp)
	if err != nil {
		return 0
	}

	return float64(t.Unix())
}

func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
