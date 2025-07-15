package constant

var ChatModeActions = map[string]string{
	"chat":    "chat",    // Requires all fields
	"restore": "restore", // Requires createdAt
	"join":    "join",    // Requires roomID
}
