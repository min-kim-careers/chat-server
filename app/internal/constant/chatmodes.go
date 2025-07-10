package constant

var ChatModes = map[string]string{
	// Flags
	"connected": "connected",
	"restored":  "restored",
	"joined":    "joined",
	"left":      "left",
	"empty":     "empty",
	"typing":    "typing",

	// Actions
	"chat":       "chat",
	"restore":    "restore", // Requires createdAt
	"join":       "join",    // Requires roomID
	"leave":      "leave",
	"disconnect": "disconnect",
}
