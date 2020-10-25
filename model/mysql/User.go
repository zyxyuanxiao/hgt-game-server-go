package model

import (
	"time"
)

// 用户主表
type User struct {
	UserId    int64     // 用户id
	Username  string    // 用户名
	Password  string    // 密码
	Avatar    string    // 头像
	Email     string    // 邮箱
	Mobile    string    // 手机号
	Gender    int       // 性别 female male
	Role      uint8     // 角色 0：普通用户 1：管理员 2：超级管理员
	Status    uint8     // 状态 0：正常 1：待定
	CreatedAt int       `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
  