package messagein

type MessageInChat struct {
	Mode    string `json:"mode"`
	TempID  string `json:"tempId"`
	Content string `json:"content"`
}

func (*MessageInChat) isMessageIn() {}
