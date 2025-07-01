package chat

import (
	"chat-server/internal/deps"
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	deps   *deps.Container
	ctx    context.Context
	cancel context.CancelFunc
	conn   *websocket.Conn
	lock   sync.Mutex
}

type Hub struct {
	Deps             *deps.Container
	Rooms            map[uuid.UUID]*Room
	RoomRegister     chan *Room
	RoomUnregister   chan *Room
	Clients          map[string]*Client
	ClientRegister   chan *Client
	ClientUnregister chan *Client
}

type Room struct {
	Hub        *Hub
	ID         uuid.UUID
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}
