package messagein

type MessageInEvent struct {
	Mode     string `json:"mode"`
	ClientID string `json:"clientId,omitempty"`
}

func (*MessageInEvent) isMessageIn() {}
