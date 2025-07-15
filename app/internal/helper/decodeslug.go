package helper

import (
	"encoding/base64"
)

func DecodeSlug(slug string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(slug)
}
