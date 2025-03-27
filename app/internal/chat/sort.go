package chat

import "time"

var ISOTimestampLayout = "2006-01-02T15:04:05.999Z"

func partition(arr []*Message, low, high int) ([]*Message, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		t1, _ := time.Parse(ISOTimestampLayout, string(arr[j].Timestamp))
		t2, _ := time.Parse(ISOTimestampLayout, string(pivot.Timestamp))
		if t1.Before(t2) {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []*Message, low, high int) []*Message {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func QuickSortStart(arr []*Message) []*Message {
	return quickSort(arr, 0, len(arr)-1)
}
