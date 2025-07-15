package messageout

type MessageOutEvent struct {
	Mode string `json:"mode"`
}

func (*MessageOutEvent) isMessageOut() {}
