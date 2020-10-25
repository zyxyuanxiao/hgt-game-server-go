package api

import (
	"encoding/json"
	"fmt"
)

type GameApi struct{}


type GameStartMessage struct {
	Args struct{
		RoomId string
	}
}
// 游戏开局
func (c *WebsocketClient) GameStart(message *WebsocketReadMessage) {
	var args *GameStartMessage
	json.Unmarshal([]byte(message.Message), &args)
}

type GamePrepareMessage struct {
	Args struct{
		RoomId string
	}
}
// 游戏准备
func (c *WebsocketClient) GamePrepare(message WebsocketReadMessage) {
	var args *GamePrepareMessage
	json.Unmarshal([]byte(message.Message), &args)
	// 找房间
	if room, ok := GroupManage[args.Args.RoomId]; ok {
		if room.Status != 0 {
			c.send <- []byte("当前游戏已经不能准备")
			return
		}
		c.GameGroupInfo.Status = 1
		room.WebsocketClient[c.UserDTO.UserId] = *c
		fmt.Println(room.WebsocketClient[c.UserDTO.UserId].GameGroupInfo.Status)
		c.send <- []byte("已经准备")
	} else {
		c.send <- []byte("没有找到房间")
	}
}
