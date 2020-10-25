package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xlstudio/wxbizdatacrypt"
	"server/app"
	"server/dto"
	model "server/model/mysql"
	"server/util"
)

type EncryptedDataUserInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
}

type AuthApi struct{}


// 小程序登录
func (api AuthApi) AppletLogin(c *gin.Context) {
	var code = app.PostParam(c, "code")
	var appId = app.PostParam(c, "appId")
	_, sessionKey := util.GetSessionKeyByCode(code.(string), appId.(string))
	var iv = app.PostParam(c, "iv")
	var encryptedData = app.PostParam(c, "encryptedData")
	pc := wxbizdatacrypt.WxBizDataCrypt{AppId: appId.(string), SessionKey: sessionKey}
	result, _ := pc.Decrypt(encryptedData.(string), iv.(string), true)
	var userInfo *EncryptedDataUserInfo
	json.Unmarshal([]byte(result.(string)), &userInfo)
	// 获取用户是否存在
	userBindModel := &model.UserBind{OpenId: userInfo.OpenID, Status: 1}
	hasBind, _ := app.DB.Get(userBindModel)
	var userId int64
	var userModel *model.User
	if hasBind {
		// 已经存在绑定用户
		userId = userBindModel.UserId
		userModel = &model.User{UserId: userId}
		app.DB.Get(userModel)
	} else {
		// 新用户
		userModel = &model.User{
			Username: userInfo.NickName,
			Avatar:   userInfo.AvatarURL,
			Gender:   userInfo.Gender,
			Role:     0,
			Status:   0,
		}
		userId, _ = app.DB.InsertOne(userModel)
		// 插入userbind
		newUserBind := model.UserBind{
			UserId:   userId,
			Nickname: userInfo.NickName,
			AppId:    appId.(string),
			OpenId:   userInfo.OpenID,
			UnionId:  userInfo.UnionID,
			Remark:   result.(string),
			Status:   1,
		}
		app.DB.InsertOne(newUserBind)
	}
	accessToken := util.BuildToken(dto.UserDTO{
		UserId:   userId,
		Username: userModel.Username,
		Role:     userModel.Role,
	})

	c.Set("data", map[string]interface{}{
		"accessToken": accessToken,
		"expire":      7 * 86400,
		"userInfo": map[string]interface{}{
			"userId":   userModel.UserId,
			"username": userModel.Username,
			"role":     userModel.Role,
		},
	})
}

