package router

import (
	"github.com/gin-gonic/gin"
	"server/websocket"
)

func Websocket(g *gin.Engine) {
	g.GET("/wss", websocket.Server)
}