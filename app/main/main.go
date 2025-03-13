package main

import (
	"chat-go/main/db"
	"context"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	// Initialise Database
	conn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		return
	}
	defer conn.Close(context.Background())
	log.Println("DB init successful")

	// Initialise Websocket
	flag.Parse()

	hub := NewHub()
	log.Println("Hub created.")

	go hub.Run()
	log.Println("Hub running.")

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		HandleWebsocketConnection(w, r, hub)
	})

	log.Printf("Websocket server starting on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
