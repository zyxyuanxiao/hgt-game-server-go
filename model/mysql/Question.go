package model

import (
	"time"
)

// 题库表
type Question struct {
	QuestionId int64     // 题库id
	Title      string    // 标题
	Desc       string    // 汤面 描述
	Content    string    // 汤底
	Sort       int64     // 排序值 值越大越前
	Hot        int64     // 热度（玩过的次数）
	Recommend  uint8     // 是否被推荐
	Status     uint8     // 状态 0：未上架 1：上架 2：删除
	CreatedAt  int64     `xorm:"created"`
	UpdatedAt  time.Time `xorm:"updated"`
}
