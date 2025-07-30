package messagein

import (
	"chat-server/internal/constants"
	"errors"
)

func validateMessageIn(m *MessageInBase) error {
	_, validMode := constants.ChatModeMap[m.Mode]

	if !validMode {
		return errors.New("invalid message mode")
	}

	return nil
}
