package api

import (
	"encoding/json"
	"fmt"
	_ "fmt"
	"server/app"
)

type GameGroupApi struct{}

type GroupInfo struct {
	RoomId          string
	Password        string
	QuestionId      int64
	GroupUserId     int64
	McUserId        int64
	Number          uint8
	Status          uint8
	WebsocketClient map[int64]WebsocketClient
}

// 群管理
var GroupManage = make(map[string]GroupInfo)

// 结束的群管理
var FinshGroupManage = make(map[string]GroupInfo)

type CreateGroupMessage struct {
	Args struct{
		Number   uint8
		Password string
	}
}
// 创房组
func (c *WebsocketClient) CreateGroup(message *WebsocketReadMessage) {
	var args *CreateGroupMessage
	json.Unmarshal([]byte(message.Message), &args)
	roomId := app.RandStr(8)
	var websocketClients = make(map[int64]WebsocketClient)
	websocketClients[c.UserDTO.UserId] = *c
	GroupManage[roomId] = GroupInfo{
		RoomId:          roomId,
		Password:        args.Args.Password,
		QuestionId:      0,
		GroupUserId:     c.UserDTO.UserId,
		McUserId:        0,
		Number:          args.Args.Number,
		Status:          0,
		WebsocketClient: websocketClients,
	}
	if c.RoomId != "" {
		// 表明之前群主不是他 @todo
	}
	fmt.Println(c.UserDTO.Username + "创建了房间id：" + roomId)
	// 设置当前群id
	c.RoomId = roomId
	c.send <- "创建成功"
}

type JoinGroupMessage struct {
	Args struct{
		RoomId string
	}
}
// 加入房间组
func (c *WebsocketClient) JoinGroup(message *WebsocketReadMessage) {
	var args *JoinGroupMessage
	json.Unmarshal([]byte(message.Message), &args)
	// 找房间
	if room, ok := GroupManage[args.Args.RoomId]; ok {
		if room.Status != 0 {
			c.error <- "当前游戏状态不能加入"
			return
		}
		room.WebsocketClient[c.UserDTO.UserId] = *c
		fmt.Println("用户：" + c.UserDTO.Username + " 加入房间："+ args.Args.RoomId)
		c.send <- "成功加入房间"
	} else {
		c.error <- "没有找到房间"
	}
}

type LeaveGroupMessage struct {
	Args struct{
		RoomId string
	}
}
// 离开群
func (c *WebsocketClient) LeaveGroup(message *WebsocketReadMessage) {
	var args *LeaveGroupMessage
	json.Unmarshal([]byte(message.Message), &args)
	// 找房间
	if room, ok := GroupManage[args.Args.RoomId]; ok {
		// 找群里面是否有这个用户
		if _, ok2 := room.WebsocketClient[c.UserDTO.UserId]; ok2 {
			delete(room.WebsocketClient, c.UserDTO.UserId)
		}
		fmt.Println("用户：" + c.UserDTO.Username + " 离开了房间："+ args.Args.RoomId)
		c.send <- "已离开房间"
	} else {
		c.error <- "没有找到房间"
	}
}

