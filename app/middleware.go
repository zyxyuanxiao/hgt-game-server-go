package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"server/exception"
	"strings"
	"time"
)

// 全局返回 json
func ReturnJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		data, _ :=  c.Get("data")
		if data != nil {
			c.JSON(http.StatusOK, gin.H{
				"_message": "success",
				"_code": 0,
				"_data": data,
			})
		}
	}
}

// 异常捕获恢复程序中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var (
					unknownErr = false // 是否未知错误
					brokenPipe bool // 断开连接
					message = "" // 响应信息
					code = 100 // 错误码
					status = 200 // 状态码
				)
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				switch err.(type) {
					case exception.LogicException:
						exception := err.(exception.LogicException)
						message = exception.Message
						code = exception.Code
						status = exception.Status
					case exception.AuthException:
						exception := err.(exception.AuthException)
						message = exception.Message
						code = exception.Code
						status = exception.Status
					default:
						unknownErr = true
				}
				go func() {
					if brokenPipe || unknownErr {
						stack := stack(3)
						httpRequest, _ := httputil.DumpRequest(c.Request, false)
						if brokenPipe {
							Logger().Error(fmt.Sprintf("%s\n%s", err, string(httpRequest)))
						} else if gin.IsDebugging() {
							headers := strings.Split(string(httpRequest), "\r\n")
							for idx, header := range headers {
								current := strings.Split(header, ":")
								if current[0] == "Authorization" {
									headers[idx] = current[0] + ": *"
								}
							}
							headersJson, _ := json.Marshal(headers)
							Logger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s",
								time.Now(), headersJson, err, stack))
						} else {
							Logger().Error(fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s",
								time.Now(), err, stack))
						}
					}
				}()
				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					if unknownErr {
						c.AbortWithStatus(http.StatusInternalServerError)
					} else {
						c.JSON(status, gin.H{
							"_message": message,
							"_code": code,
							"_data": nil,
						})
					}
				}

			}
		}()
		c.Next()
	}
}

// 请求日志
func LoggerToFile() gin.HandlerFunc {
	logger := Logger()
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		//日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}