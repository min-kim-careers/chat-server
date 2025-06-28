package main

import (
	"chat-server/internal/api"
	"chat-server/internal/cache"
	"chat-server/internal/chat"
	"chat-server/internal/db"
	"chat-server/internal/websocket"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Database
	newDB := db.NewDB()
	newDB.CreateMessagesTable()
	log.Println("DB init successful.")

	// Cache
	newCache := cache.NewCache()
	log.Println("Cache init successful.")

	// Hub
	newHub := chat.NewHub(newDB, newCache)
	log.Println("Hub created.")

	newHub.Run()
	log.Println("Hub running.")

	apiAddr := flag.String("apiAddr", ":8081", "API server address")
	wsAddr := flag.String("wsAddr", ":8080", "WebSocket server address")
	flag.Parse()

	// API server
	go func() {
		router := gin.Default()
		api.RegisterMessageRoutes(router.Group("/messages"), newCache, newDB)
		log.Println("Chat API server running on", *apiAddr)
		if err := router.Run(*apiAddr); err != nil {
			panic(err)
		}
	}()

	// WS server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket.WebsocketHandler(w, r, newHub)
	})
	log.Println("Chat WS server running on", *wsAddr)
	if err := http.ListenAndServe(*wsAddr, nil); err != nil {
		panic(err)
	}
}
