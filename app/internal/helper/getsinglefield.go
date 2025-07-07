package helper

import (
	"log"

	"github.com/buger/jsonparser"
)

func GetSingleField(p []byte, field ...string) []byte {
	val, _, _, err := jsonparser.Get(p, field...)
	if err != nil {
		log.Printf("Error getting single field: %v", err)
		return nil
	}
	return val
}
