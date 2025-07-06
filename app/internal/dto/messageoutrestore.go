package dto

type MessageOutRestore struct {
	Mode     string            `json:"mode"`
	Messages []*MessageOutChat `json:"messages"`
}
