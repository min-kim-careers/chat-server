package messageout

type MessageOutBase struct {
	Mode string `json:"mode"`
}

func (*MessageOutBase) isMessageOut() {}
