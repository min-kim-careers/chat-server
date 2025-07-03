package helper

import (
	"chat-server/internal/dto"
)

func partition(arr []*dto.Message, low, high int) ([]*dto.Message, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		t1 := arr[j].CreatedAt
		t2 := pivot.CreatedAt
		if t1.Before(t2) {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []*dto.Message, low, high int) []*dto.Message {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func QuickSortStart(msgs []*dto.Message) []*dto.Message {
	return quickSort(msgs, 0, len(msgs)-1)
}
