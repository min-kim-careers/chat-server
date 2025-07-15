package messagein

type MessageInBase struct {
	Mode string `json:"mode"`
}

func (*MessageInBase) isMessageIn() {}
