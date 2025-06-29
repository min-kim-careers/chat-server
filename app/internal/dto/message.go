package dto

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

var validMessageTypes = map[string]bool{
	"connect":    true,
	"disconnect": true,
	"chat":       true,
	"restore":    true,
	"empty":      true,
}

type Message struct {
	ID          int             `json:"id"`
	MessageType string          `json:"messageType"`
	RoomID      string          `json:"roomId"`
	ClientID    string          `json:"clientId"`
	CreatedAt   time.Time       `json:"createdAt"`
	Data        json.RawMessage `json:"data"`
}

func validateMessage(msg *Message) bool {
	validate := validator.New()

	err := validate.Struct(msg)
	if err != nil {
		log.Println("Invalid message:", err)
		return false
	}

	_, valid := validMessageTypes[msg.MessageType]
	if !valid {
		log.Println("Invalid message type:", msg.MessageType)
		return false
	}

	return true
}

func DeserializeMessage(jsonData []byte) *Message {
	var msg Message

	err := json.Unmarshal(jsonData, &msg)
	if err != nil {
		log.Println("Error deserializing message:", string(jsonData))
		log.Println(err)
		return nil
	}

	if !validateMessage(&msg) {
		return nil
	}

	return &msg
}

func SerializeMessage(m *Message) []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println("Error serializing message:", err)
		return nil
	}

	return data
}

// func PrintMessage(m *MessageDTO) {
// 	res, err := json.MarshalIndent(m, "", "  ")
// 	if err != nil {
// 		log.Println("Error printing message:", err)
// 		return
// 	}
// 	log.Println(res)
// }

// func PrintJson(j []byte) {
// 	var buffer bytes.Buffer
// 	err := json.Indent(&buffer, j, "", "\t")
// 	if err != nil {
// 		log.Println("Error prettifying JSON:", j)
// 		return
// 	}
// 	log.Println(buffer.String())
// }
