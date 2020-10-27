package websocket

var ProtocolActions = make(map[int64]func(struct{}, *Client, *ReadMessage))
var ProtocolStruct = make(map[int64]struct{})
var ActionProtocols = make(map[string]int64)

// 加载协议
func LoadProtocol() {
	ProtocolActions[2001] = (*GameRoom).CreateGroup
	ProtocolStruct[2001] = &GameRoom{}
}