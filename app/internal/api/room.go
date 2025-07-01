package api

import (
	"chat-server/internal/deps"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterRoomRoutes(rg *gin.RouterGroup, d *deps.Container) {
	rg.POST("/rooms/register", func(c *gin.Context) {
		registerRoom(c, d)
	})
}

type RegisterRoomRequest struct {
	ItemID  string    `json:"itemId"`
	Client1 uuid.UUID `json:"client1"`
	Client2 uuid.UUID `json:"client2"`
}

type RegisterRoomResponse struct {
	RoomID string `json:"roomId"`
}

func registerRoom(c *gin.Context, d *deps.Container) {
	var req RegisterRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(req)
		c.JSON(http.StatusBadRequest, APIError{Message: "Bad request"})
		return
	}

	room, err := d.Services.Room.RegisterRoom(c.Request.Context(), req.ItemID, req.Client1, req.Client2)
	if err != nil {
		log.Printf("Failed to register room for item <%s>: %v", req.ItemID, err)
		c.JSON(http.StatusInternalServerError, APIError{Message: "Failed to register room"})
	}

	c.JSON(http.StatusCreated, APIResponse{Data: RegisterRoomResponse{RoomID: room.ID.String()}})
}
