package main

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
	MessageType    MessageType    `json:"message_type"`
	RoomID         RoomID         `json:"room_id"`
	ClientID       ClientID       `json:"client_id"`
	MessageContent MessageContent `json:"message_content"`
}

func NewMessage(messageType MessageType, roomID RoomID, clientID ClientID, messageContent MessageContent) *Message {
	m := Message{}

	m.MessageType = messageType
	m.RoomID = roomID
	m.ClientID = clientID
	m.MessageContent = messageContent

	return &m
}

func Deserialize(jsonData []byte) *Message {
	var m Message

	err := json.Unmarshal(jsonData, &m)
	if err != nil {
		log.Println("Error deserializing message")
		return nil
	}

	_, valid := validMessageTypes[m.MessageType]
	if !valid {
		log.Println("Invalid message type:", m.MessageType)
		return nil
	}

	return &m
}

func Serialize(m *Message) []byte {
	data, err := json.Marshal(m)
	if err != nil {
		log.Println("Error serializing message")
		return nil
	}

	return data
}

func printMessage(m *Message) {
	res, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println("Error printing message:", err)
		return
	}
	log.Println(res)
}

func prettyJson(j []byte) {
	var buffer bytes.Buffer
	err := json.Indent(&buffer, j, "", "\t")
	if err != nil {
		log.Println("Error prettifying JSON:", j)
		return
	}
	log.Println(buffer)
}
