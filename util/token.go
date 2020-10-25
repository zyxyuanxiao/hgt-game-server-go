package util

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"server/app"
	"server/dto"
	"server/exception"
)

// 获取token
func BuildToken(userDto dto.UserDTO) string {
	client := app.Redis.GetRedisByPool().Get()
	defer client.Close()
	token := uuid.NewV4().String()
	userByte, _ := json.Marshal(userDto)
	userJson := fmt.Sprintf("%s", userByte)
	client.Do("set", token, userJson, "ex", 7*86400)

	return token
}

// 校验token
func CheckToken(token string) dto.UserDTO {
	client := app.Redis.GetRedisByPool().Get()
	defer client.Close()
	var userDto dto.UserDTO
	userJson, _ := redis.String(client.Do("get", token))
	if userJson == "" {
		// token已经过期
		exception.Auth("token expired")
	}
	// 自动续期
	client.Do("set", token, userJson, "ex", 7*86400)
	json.Unmarshal([]byte(userJson), &userDto)

	return userDto
}
