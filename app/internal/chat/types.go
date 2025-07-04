package chat

import (
	"chat-server/internal/service"
	"context"

	"github.com/gorilla/websocket"
)

type Hub struct {
	svc              *service.Services
	rooms            map[string]*Room
	roomRegister     chan *Room
	roomUnregister   chan *Room
	clients          map[string]*Client
	clientRegister   chan *Client
	clientUnregister chan *Client
}

type Room struct {
	hub              *Hub
	id               string
	clients          map[string]*Client
	clientRegister   chan *Client
	clientUnregister chan *Client
}

type Client struct {
	hub     *Hub
	room    *Room
	id      string
	ctx     context.Context
	cancel  context.CancelFunc
	conn    *websocket.Conn
	channel chan []byte
}
