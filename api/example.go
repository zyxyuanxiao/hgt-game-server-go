package api

import (
	"github.com/gin-gonic/gin"
	"server/app"
	"server/exception"
	model "server/model/mysql"
	"strconv"
)

type ExampleApi struct {}

func (u ExampleApi) Test(c *gin.Context)  {
	userId, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	var user = &model.User{UserId: userId}
	has, _ := app.DB.Get(user)
	if !has {
		exception.Logic("用户不存在")
	}

	var data = map[string]string{
		"username": user.Username,
	}

	c.Set("data", data)
}
