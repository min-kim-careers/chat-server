package main

import (
	"chat-server/internal/cache"
	"chat-server/internal/chat"
	"chat-server/internal/db"
	"chat-server/internal/ws"
	"flag"
	"log"
	"net/http"
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

	// Ports
	addr := flag.String("addr", ":8080", "WebSocket server address")
	flag.Parse()

	// WS server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketHandler(w, r, newHub)
	})
	log.Println("WS server running on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
