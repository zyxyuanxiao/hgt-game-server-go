package model

import (
	"time"
)

// 用户绑定表
type UserBind struct {
	UserBindId int64     // 用户授权绑定id
	UserId     int64     // 用户id
	Category   uint8     // 类型 1：微信
	Nickname   string    // 第三方名称
	AppId      string    // 第三方appId
	OpenId     string    // 第三方账号
	UnionId    string    // 第三方管理唯一id
	Remark     string    // 备注
	Status     uint8     // 状态 0：未绑定 1：绑定中
	CreatedAt  int       `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
}
