package chat

import (
	"chat-server/internal/service"
	"context"

	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	svc    *service.Services
	ctx    context.Context
	cancel context.CancelFunc
	conn   *websocket.Conn
}

type Hub struct {
	Svc              *service.Services
	Rooms            map[string]*Room
	RoomRegister     chan *Room
	RoomUnregister   chan *Room
	Clients          map[string]*Client
	ClientRegister   chan *Client
	ClientUnregister chan *Client
}

type Room struct {
	Hub        *Hub
	ID         string
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}
