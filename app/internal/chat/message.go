package chat

import (
	"bytes"
	"encoding/json"
	"log"
)

type (
	MessageType    string
	ItemID         string
	TargetID       string
	MessageContent string
	Timestamp      string
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
}

type Message struct {
	Type      MessageType    `json:"message_type"`
	RoomID    RoomID         `json:"room_id"`
	ClientID  ClientID       `json:"client_id"`
	Content   MessageContent `json:"message_content"`
	Timestamp Timestamp      `json:"timestamp"`
}

func NewMessage(messageType MessageType, roomID RoomID, clientID ClientID, messageContent MessageContent, timestamp Timestamp) *Message {
	m := Message{}

	m.Type = messageType
	m.RoomID = roomID
	m.ClientID = clientID
	m.Content = messageContent
	m.Timestamp = timestamp

	return &m
}

func DeserializeMessage(jsonData []byte) *Message {
	var msg Message

	err := json.Unmarshal(jsonData, &msg)
	if err != nil {
		log.Println("Error deserializing message")
		return nil
	}

	_, valid := validMessageTypes[msg.Type]
	if !valid {
		log.Println("Invalid message type:", msg.Type)
		return nil
	}

	return &msg
}

func SerializeMessage(m *Message) []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println("Error serializing message")
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
	log.Println(buffer)
}
