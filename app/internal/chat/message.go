package chat

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

type (
	MessageType    string
	ItemID         string
	TargetID       string
	Timestamp      string
	MessageContent any
)

var validMessageTypes = map[MessageType]bool{
	"connect":    true,
	"disconnect": true,
	"enter":      true,
	"leave":      true,
	"chat":       true,
	"typing":     true,
	"edit":       true,
	"delete":     true,
	"restore":    true,
}

type Message struct {
	Type      MessageType    `json:"message_type" validate:"required"`
	RoomID    RoomID         `json:"room_id" validate:"required"`
	ClientID  ClientID       `json:"client_id" validate:"required"`
	Timestamp Timestamp      `json:"timestamp"`
	Content   MessageContent `json:"message_content"`
}

func validateMessage(msg *Message) bool {
	validate := validator.New()

	err := validate.Struct(msg)
	if err != nil {
		log.Println("Invalid message:", err)
		return false
	}

	_, valid := validMessageTypes[msg.Type]
	if !valid {
		log.Println("Invalid message type:", msg.Type)
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

func PrintMessage(m *Message) {
	res, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println("Error printing message:", err)
		return
	}
	log.Println(res)
}

func PrintJson(j []byte) {
	var buffer bytes.Buffer
	err := json.Indent(&buffer, j, "", "\t")
	if err != nil {
		log.Println("Error prettifying JSON:", j)
		return
	}
	log.Println(buffer.String())
}
