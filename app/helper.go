package app

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"time"
)

func PostParam(c *gin.Context, key string) interface{} {
	value := c.Request.FormValue(key)
	if value == "" {
		// 兼容rawBody
		return RowBody(c, key)
	}

	return value
}

func RowBody(c *gin.Context, key string) interface{} {
	defer c.Request.Body.Close()
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	// 把刚刚读出来的再写进去
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// 解析参数 存入map
	var rawBodyParams map[string]string
	// 用于存放参数key=value数据
	json.Unmarshal(bodyBytes, &rawBodyParams)
	if value, ok := rawBodyParams[key]; ok{
		return value
	}

	return ""
}

// 随机生成字符串
func RandStr(length int) string {
	str := "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}