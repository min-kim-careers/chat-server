package main

import (
	"encoding/json"
	"log"
)

type Message struct {
	Type     string `json:"type"`      // Message type (e.g., "enter", "talk")
	ItemID   string `json:"item_id"`   // Identifier for the chat room
	SenderID string `json:"sender_id"` // Sender's identifier
	TargetID string `json:"target_id"`
	Message  string `json:"message"` // The content of the message
}

func Deserialize(jsonData []byte) *Message {
	var message Message
	err := json.Unmarshal(jsonData, &message)
	if err != nil {
		log.Println("Failed to deserialize message")
		return nil
	}
	return &message
}

func Serialize(m *Message) []byte {
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Println("Failed to serialize message")
		return nil
	}
	return jsonData
}
