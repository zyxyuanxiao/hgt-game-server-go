package main

import (
	"server/app"
	"server/router"
)

func main() {
	// 框架初始化
	g := app.Init()
	// 初始化路由
	router.Init(g)

	// 监听服务
	g.Run(":8899")
}
