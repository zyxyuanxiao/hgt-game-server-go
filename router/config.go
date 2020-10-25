package router

import (
	"github.com/gin-gonic/gin"
)

func Init(g *gin.Engine) {
	loadRouter(g)
}

// 加载路由
func loadRouter(g *gin.Engine) {
	// 示例
	Example(g)
	// 认证
	Auth(g)
	// websocket
	Websocket(g)
}
