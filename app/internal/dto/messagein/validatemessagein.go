package messagein

import (
	"chat-server/internal/constant"
	"fmt"
)

func validateMessageIn(m *MessageInBase) (bool, error) {
	_, isEvent := constant.ChatModeEvents[m.Mode]

	_, isAction := constant.ChatModeActions[m.Mode]

	validMode := isEvent || isAction
	if !validMode {
		return false, fmt.Errorf("invalid message mode: %s", m.Mode)
	}

	return isEvent, nil
}
