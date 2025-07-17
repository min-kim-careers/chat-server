package messagein

import (
	"chat-server/internal/constants"
	"fmt"
)

func validateMessageIn(m *MessageInBase) (bool, error) {
	_, isEvent := constants.ChatModeEvents[m.Mode]

	_, isAction := constants.ChatModeActions[m.Mode]

	validMode := isEvent || isAction
	if !validMode {
		return false, fmt.Errorf("invalid message mode: %s", m.Mode)
	}

	return isEvent, nil
}
