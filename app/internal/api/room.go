package api

import (
	"chat-server/internal/deps"

	"github.com/gin-gonic/gin"
)

func RegisterRoomRoutes(rg *gin.RouterGroup, deps *deps.Container) {
	rg.POST("/", func(c *gin.Context) {
		addRoom(c, deps)
	})
}

func addRoom(c *gin.Context, deps *deps.Container) {

}
