package main

import (
	"chat-server/internal/api"
	"chat-server/internal/chat"
	"chat-server/internal/handler"
	"chat-server/internal/service"
	"flag"
	"log"
	"net/http"

	grmon "github.com/bcicen/grmon/agent"
	"github.com/gin-gonic/gin"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	grmon.Start()

	services := service.NewServices()

	// Hub
	hub := chat.NewHub(services)
	log.Println("Hub created.")

	hub.Run()
	log.Println("Hub running.")

	apiAddr := flag.String("apiAddr", ":8081", "API server address")
	wsAddr := flag.String("wsAddr", ":8080", "WebSocket server address")
	flag.Parse()

	// API server
	go func() {
		router := gin.Default()
		api.RegisterMessageRoutes(router.Group(""), services)
		api.RegisterRoomRoutes(router.Group(""), services)
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
