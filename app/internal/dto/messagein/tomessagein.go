package messagein

import (
	"chat-server/internal/helper"
	"encoding/json"
	"errors"
	"fmt"
)

func ToMessageIn(p []byte) (MessageIn, error) {
	var b MessageInBase
	if err := json.Unmarshal(p, &b); err != nil {
		return nil, err
	}

	isEvent, err := validateMessageIn(&b)
	if err != nil {
		return nil, err
	}

	if isEvent {
		var e MessageInEvent
		if err := json.Unmarshal(p, &e); err != nil {
			return nil, err
		}
		return &e, nil
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

	case "restore":
		var r MessageInRestore

		if err := json.Unmarshal(p, &r); err != nil {
			return nil, err
		}
		if r.CreatedAt.IsZero() {
			return nil, errors.New("restore missing createdAt")
		}
		return &r, nil
	}

	return nil, fmt.Errorf("unknown type %s", b.Mode)
}
