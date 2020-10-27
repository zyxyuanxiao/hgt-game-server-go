package api

import (
	"encoding/json"
	"server/websocket"
)

type GameApi struct{}


type GameStartMessage struct {
	Args struct{
		RoomId string
	}
}
// 游戏开局
func GameStart(c *websocket.Client, message *websocket.ReadMessage) {
	var args *GameStartMessage
	json.Unmarshal([]byte(message.Message), &args)
}

type GamePrepareMessage struct {
	Args struct{
		RoomId string
	}
}
// 游戏准备
func GamePrepare(c *websocket.Client, message websocket.ReadMessage) {
	var args *GamePrepareMessage
	json.Unmarshal([]byte(message.Message), &args)
	// 找房间
	if room, ok := GroupManage[args.Args.RoomId]; ok {
		if room.Status != 0 {
			c.Send <- []byte("当前游戏已经不能准备")
			return
		}
		//c.GameGroupInfo.Status = 1
		room.WebsocketClient[c.UserDTO.UserId] = *c
		//fmt.Println(room.WebsocketClient[c.UserDTO.UserId].GameGroupInfo.Status)
		c.Send <- []byte("已经准备")
	} else {
		c.Send <- []byte("没有找到房间")
	}
}
