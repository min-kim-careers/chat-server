package helper

import (
	"encoding/base64"
)

func EncodeSlug(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)

}
