package websocket

import (
	"github.com/golang/protobuf/proto"
)

var ProtocolActions = make(map[int64]func(*proto.Message, *Client, *ReadMessage))
var ProtocolStruct = make(map[int64]proto.Message)
var ActionProtocols = make(map[string]int64)

// 加载协议
func LoadProtocol() {
	ProtocolActions[2001] = (*CreateRoomReq).CreateGroup
	ProtocolStruct[2001] = &CreateRoomReq{}
}