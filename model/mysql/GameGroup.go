package model

import (
	"time"
)

// 游戏群表
type GameGroup struct {
	GameId    int64     // 游戏id
	UserId    string    // 用户id
	Status    uint8     // 状态 0：未准备 1：准备中 2：游戏中 3：掉线 4：被踢（不能再进）
	CreatedAt int64     `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
