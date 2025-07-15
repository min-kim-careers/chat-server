package messageout

type MessageOutRestored struct {
	Mode     string            `json:"mode"`
	Messages []*MessageOutChat `json:"messages"`
}

func (*MessageOutRestored) isMessageOut() {}
