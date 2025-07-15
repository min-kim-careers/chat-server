package messageout

import (
	"encoding/json"
)

func ToRawMessageOut(m MessageOut) ([]byte, error) {
	p, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return p, nil
}
