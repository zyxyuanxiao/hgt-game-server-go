package app

import (
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

// 项目目录
var Path, _ = os.Getwd()

// 环境变量
var ENV = os.Getenv("env")

func Init() *gin.Engine {
	// 设置为发布模式
	if ENV == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	// 加入错误捕获 全局返回json 中间件
	g.Use(LoggerToFile(), Recovery(), ReturnJson())
	// 加载配置
	LoadConfig()
	// 加载数据库
	LoadDB()
	// 加载redis
	LoadRedis()
	// 支持 cron
	go CronStart()

	return g
}