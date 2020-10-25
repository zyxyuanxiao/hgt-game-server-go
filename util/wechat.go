package util

import (
	"github.com/medivhzhan/weapp/v2"
	"server/app"
)

func GetSessionKeyByCode(code string, appId string) (string, string){
	res, err := weapp.Login(appId, app.AppletConfig[appId].Secret, code)

	if err != nil {
		// 处理一般错误信息
		panic(err)
	}

	if err := res.GetResponseError(); err !=nil {
		// 处理微信返回错误信息
		panic(err)
	}
	return res.OpenID, res.SessionKey
}
