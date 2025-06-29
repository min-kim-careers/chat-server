package main

import (
	"chat-server/internal/api"
	"chat-server/internal/chat"
	"chat-server/internal/deps"
	"chat-server/internal/handler"
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Dependencies
	ctx := context.Background()
	deps := deps.NewContainer(ctx)

	// Hub
	hub := chat.NewHub(deps)
	log.Println("Hub created.")

	hub.Run()
	log.Println("Hub running.")

	apiAddr := flag.String("apiAddr", ":8081", "API server address")
	wsAddr := flag.String("wsAddr", ":8080", "WebSocket server address")
	flag.Parse()

	// API server
	go func() {
		router := gin.Default()
		api.RegisterMessageRoutes(router.Group("/message"), deps)
		api.RegisterRoomRoutes(router.Group("/room"), deps)
		log.Println("Chat API server running on", *apiAddr)
		if err := router.Run(*apiAddr); err != nil {
			panic(err)
		}
	}()

	// WS server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.WebsocketHandler(w, r, hub)
	})
	log.Println("Chat WS server running on", *wsAddr)
	if err := http.ListenAndServe(*wsAddr, nil); err != nil {
		panic(err)
	}
}
