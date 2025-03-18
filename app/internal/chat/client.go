package chat

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientID string

type Client struct {
	id     ClientID
	wsConn *websocket.Conn
	cache  *Cache
}

func NewClient(id ClientID, roomID RoomID, wsConn *websocket.Conn, cache *Cache) *Client {
	return &Client{
		id:     id,
		wsConn: wsConn,
		cache:  cache,
	}
}

// Send to client
func (client *Client) Send(room *Room) {
	streams := client.cache.GetReadStreams(string(room.id))
	if streams == nil {
		return
	}

	for _, stream := range streams {
		for _, msg := range stream.Messages {
			msgInterface, ok := msg.Values["message"]
			if !ok {
				log.Println("No message in stream message:", msg)
				continue
			}

			msgJson := []byte(msgInterface.(string))

			msg := DeserializeMessage(msgJson)
			if msg == nil {
				continue
			}

			if msg.ClientID == client.id {
				continue
			}

			err := client.wsConn.WriteMessage(websocket.TextMessage, msgJson)
			if err != nil {
				log.Printf("Error sending message to client <%s>: %v", client.id, err)
			} else {
				log.Printf("Sent message to client <%s>: %v", client.id, msg)
			}
		}
	}
}

// Read from client
func (client *Client) Read(room *Room) {
	for {
		_, msgJson, err := client.wsConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client <%s>: %v", client.id, err)
			continue
		}

		client.cache.AddToStream(string(room.id), string(msgJson))
	}
}

func (client *Client) Run(room *Room) {
	go client.Read(room)
	go client.Send(room)
}
