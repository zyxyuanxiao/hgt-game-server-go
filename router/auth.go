package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)
var authApi = api.AuthApi{}

func Auth(g *gin.Engine) {
	g.POST("/auth/appletLogin", authApi.AppletLogin)
}
