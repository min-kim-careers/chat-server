package helper

import (
	"log"

	"github.com/buger/jsonparser"
)

func GetFieldValue(p []byte, field ...string) []byte {
	val, _, _, err := jsonparser.Get(p, field...)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	return val
}
