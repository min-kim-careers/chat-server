package messageout

import (
	"encoding/json"
	"log"
)

func ToRawMessageOut(m MessageOut) ([]byte, error) {
	p, err := json.Marshal(m)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	return p, nil
}
