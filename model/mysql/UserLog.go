package model

import (
	"time"
)

// 用户日志表
type UserLog struct {
	UserLogId    int64     // 用户日志操作表id
	UserId       int64     // 用户id
	HandleUserId int64     // 操作的用户id
	Remark       string    // 备注
	CreatedAt    int64     `xorm:"created"`
	UpdatedAt    time.Time `xorm:"updated"`
}
