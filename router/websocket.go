package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)
var websocketApi = api.WebsocketApi{}

func Websocket(g *gin.Engine) {
	g.GET("/wss", websocketApi.Websocket)
}
