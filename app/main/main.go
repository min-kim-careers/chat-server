package main

import (
	"chat-go/main/db"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {

	// Initialise Database
	conn, err := db.InitDB()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())
	fmt.Println("DB init successful")

	// Initialise Websocket
	flag.Parse()

	hub := NewHub(make(map[string]*Room), make(map[string]*Client))
	go hub.Run()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		HandleWebsocketConnection(w, r, hub)
	})

	log.Printf("Websocket server starting on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

}
