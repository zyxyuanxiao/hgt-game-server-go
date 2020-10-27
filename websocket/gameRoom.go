package websocket

import (
	"encoding/json"
	"fmt"
	"server/app"
	"server/protobuf/soup"
)

type GroupInfo struct {
	RoomId          string
	Password        string
	QuestionId      int64
	GroupUserId     int64
	McUserId        int64
	Number          uint8
	Status          uint8
	WebsocketClient map[int64]Client
}

// 群管理
var GroupManage = make(map[string]GroupInfo)

type CreateGroupMessage struct {
	Args struct{
		Number   uint8
		Password string
	}
}

type CreateRoomReq struct {
	*soup.CreateRoomReq
}

// 创房组
func (room *CreateRoomReq) CreateGroup(c *Client, message *ReadMessage) {
	var args *CreateGroupMessage
	json.Unmarshal([]byte(message.Message), &args)
	roomId := app.RandStr(8)
	var websocketClients = make(map[int64]Client)
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
	//if c.RoomId != "" {
	//	// 表明之前群主不是他 @todo
	//}
	fmt.Println(c.UserDTO.Username + "创建了房间id：" + roomId)
	// 设置当前群id
	//c.RoomId = roomId
	c.Send <- "创建成功"
}

type JoinGroupMessage struct {
	Args struct{
		RoomId string
	}
}
// 加入房间组
func JoinGroup(c *Client, message *ReadMessage) {
	var args *JoinGroupMessage
	json.Unmarshal([]byte(message.Message), &args)
	// 找房间
	if room, ok := GroupManage[args.Args.RoomId]; ok {
		if room.Status != 0 {
			c.Error <- "当前游戏状态不能加入"
			return
		}
		room.WebsocketClient[c.UserDTO.UserId] = *c
		fmt.Println("用户：" + c.UserDTO.Username + " 加入房间："+ args.Args.RoomId)
		c.Send <- "成功加入房间"
	} else {
		c.Error <- "没有找到房间"
	}
}

type LeaveGroupMessage struct {
	Args struct{
		RoomId string
	}
}
// 离开群
func LeaveGroup(c *Client, message *ReadMessage) {
	var args *LeaveGroupMessage
	json.Unmarshal([]byte(message.Message), &args)
	// 找房间
	if room, ok := GroupManage[args.Args.RoomId]; ok {
		// 找群里面是否有这个用户
		if _, ok2 := room.WebsocketClient[c.UserDTO.UserId]; ok2 {
			delete(room.WebsocketClient, c.UserDTO.UserId)
		}
		fmt.Println("用户：" + c.UserDTO.Username + " 离开了房间："+ args.Args.RoomId)
		c.Send <- "已离开房间"
	} else {
		c.Error <- "没有找到房间"
	}
}

