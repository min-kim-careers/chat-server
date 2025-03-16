package main

import (
	"chat-go/internal/chat"
	"chat-go/internal/server"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	// Initialising Database
	newDB := chat.NewDB()
	newDB.CreateMessageTable()
	defer newDB.Pool.Close()
	log.Println("DB init successful.")

	newCache := chat.NewCache()
	defer newCache.Client.Close()
	log.Println("Cache init successful.")

	// Initialising Websocket
	flag.Parse()

	newHub := chat.NewHub(newDB)
	log.Println("Hub created.")

	newHub.Run()
	log.Println("Hub running.")

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWebsocketConnection(w, r, newHub)
	})

	log.Printf("Websocket server starting on %s.", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
