package model

import (
	"time"
)

// 游戏消息表
type GameMsg struct {
	GameMsgId int64     // 游戏消息id
	UserId    string    // 用户id
	Content   string    // 消息内容
	Result    uint8     // 结果 0：无结果 1：是 2：不是 3：是或不是 4：无关
	Status    uint8     // 状态 1：正常 2：删除
	CreatedAt int64     `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
