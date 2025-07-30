package messagein

import (
	"encoding/json"
	"errors"
	"log"
)

func ToRawMessageIn(m MessageIn) ([]byte, error) {
	v, ok := m.(*MessageInBase)
	if !ok {
		return nil, errors.New("wrong message in type")
	}

	if err := validateMessageIn(v); err != nil {
		return nil, err
	}

	p, err := json.Marshal(m)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}

	return p, nil
}
