package api

import (
	"chat-server/internal/helper"
	"chat-server/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoomRoutes(rg *gin.RouterGroup, s *service.Services) {
	rg.POST("/rooms/register", func(c *gin.Context) {
		registerRoom(c, s)
	})

	rg.GET("/rooms/client/:clientID", func(c *gin.Context) {
		getRoomsByClient(c, s)
	})

	rg.DELETE("/rooms/roomSlug/:roomSlug", func(c *gin.Context) {
		deleteRoomById(c, s)
	})
}

type RegisterRoomRequest struct {
	ItemID   string    `json:"itemId"`
	ClientID uuid.UUID `json:"clientId"`
	PeerID   uuid.UUID `json:"peerId"`
}

func registerRoom(c *gin.Context, s *service.Services) {
	var req RegisterRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIError{Message: "bad request"})
		return
	}

	if req.ClientID == uuid.Nil || req.PeerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, APIError{Message: "undefined client"})
		return
	}

	room, err := s.Room.RegisterRoom(c.Request.Context(), req.ItemID, req.ClientID, req.PeerID)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, APIError{Message: "failed to register room"})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{Data: gin.H{
		"room": room,
	}})
}

func getRoomsByClient(c *gin.Context, s *service.Services) {
	clientID := c.Param("clientID")
	_clientID, err := uuid.Parse(clientID)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusBadRequest, APIError{Message: "invalid client id"})
		return
	}

	rooms, err := s.Room.GetAllRoomsByClient(c, _clientID)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, APIError{Message: "failed to get all rooms"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{Data: gin.H{
		"rooms": rooms,
	}})
}

func deleteRoomById(c *gin.Context, s *service.Services) {
	roomSlug := c.Param("roomSlug")
	roomID := helper.RoomSlugToID(roomSlug)
	if roomID == nil {
		c.JSON(http.StatusBadRequest, APIError{Message: "invalid room id"})
		return
	}

	err := s.Room.DeleteRoomById(c, *roomID)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusInternalServerError, APIError{Message: "failed to delete room by id"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{Data: gin.H{
		"message": "deleted",
	}})
}
