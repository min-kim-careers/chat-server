package main

import (
	"chat-server/internal/api"
	"chat-server/internal/cache"
	"chat-server/internal/chat"
	"chat-server/internal/db"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	// Database
	newDB := db.NewDB()
	newDB.CreateMessageTable()
	log.Println("DB init successful.")

	// Cache
	newCache := cache.NewCache()
	log.Println("Cache init successful.")

	// Hub
	newHub := chat.NewHub(newDB, newCache)
	log.Println("Hub created.")

	newHub.Run()
	log.Println("Hub running.")

	// Websocket
	flag.Parse()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		api.HandleWebsocketConnection(w, r, newHub)
	})

	log.Printf("Websocket server starting on %s.", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
