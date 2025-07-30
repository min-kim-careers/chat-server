package messagein

import (
	"chat-server/internal/helper"
	"encoding/json"
	"errors"
)

func ToMessageIn(p []byte) (MessageIn, error) {
	var b MessageInBase
	if err := json.Unmarshal(p, &b); err != nil {
		return nil, err
	}

	if err := validateMessageIn(&b); err != nil {
		return nil, err
	}

	switch b.Mode {

	case "chat":
		var c MessageInChat
		if err := json.Unmarshal(p, &c); err != nil {
			return nil, err
		}
		return &c, nil

	case "join":
		var j MessageInJoin
		if err := json.Unmarshal(p, &j); err != nil {
			return nil, err
		}
		roomID := helper.RoomSlugToID(j.RoomSlug)
		if roomID == nil {
			return nil, errors.New("invalid room slug")
		}
		j.RoomID = roomID
		return &j, nil

	default:
		var e MessageInEvent
		if err := json.Unmarshal(p, &e); err != nil {
			return nil, err
		}
		return &e, nil
	}
}
