package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"server/dto"
	"server/protobuf"
	"server/util"
	"time"
)

type ClientManager struct {
	clients    map[int64]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	UserDTO       dto.UserDTO
	socket        *websocket.Conn
	Send          chan interface{}
	Error         chan interface{}
	InsertTime    string
}

var websocketManager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[int64]*Client),
}

type ReadMessage struct {
	Message string
	Action  string
	Args    interface{}
}

// 服务
func Server(c *gin.Context) {
	//解析一个连接
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	token := c.Query("Authorization")
	if token == "" {
		conn.WriteMessage(websocket.TextMessage, []byte("Illegal connection"))
		conn.Close()
		return
	}
	go websocketManager.start()
	// 解析token
	userDto := util.CheckToken(token)

	//初始化一个客户端对象
	client := &Client{
		UserDTO:    userDto,
		socket:     conn,
		Send:       make(chan interface{}),
		Error:      make(chan interface{}),
		InsertTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	//把这个对象发送给 管道
	websocketManager.register <- client
	// 协程接收输出信息
	go client.write()
	go client.read()
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register: //新客户端加入
			manager.clients[conn.UserDTO.UserId] = conn
			fmt.Println("新用户加入:"+conn.UserDTO.Username, "加入时间："+conn.InsertTime)
			fmt.Println("当前总用户数量register：", len(manager.clients))
		case conn := <-manager.unregister: // 离开
			if _, ok := manager.clients[conn.UserDTO.UserId]; ok {
				close(conn.Send)
				delete(manager.clients, conn.UserDTO.UserId)
				fmt.Println("用户离开：" + conn.UserDTO.Username)
				fmt.Println("当前总用户数量unregister：", len(manager.clients))
			}
		case message := <-manager.broadcast: //读到广播管道数据后的处理
			fmt.Println("当前总用户数量broadcast：", len(manager.clients))
			for _, conn := range manager.clients {
				select {
				case conn.Send <- message: //调用发送给全体客户端
				default:
					// 重新上来之后挤掉了 @todo
					// 关闭连接
					close(conn.Send)
					delete(manager.clients, conn.UserDTO.UserId)
				}
			}
		}
	}
}

// 广播数据 除了自己
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for _, conn := range manager.clients {
		if conn != ignore {
			conn.Send <- message //发送的数据写入所有的 websocket 连接 管道
		}
	}
}

// 写入管道后激活这个进程
func (c *Client) write() {
	defer func() {
		// 程序退出 关闭链接
		websocketManager.unregister <- c
		c.socket.Close()
	}()

	for {
		select {
		case data, ok := <-c.Send: //这个管道有了数据 写这个消息出去
			if !ok {
				// 发送关闭提示
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			res, _ := json.Marshal(map[string]interface{}{
				"_code":    0,
				"_message": "success",
				"_data":    data,
			})
			err := c.socket.WriteMessage(websocket.TextMessage, res)
			if err != nil {
				// 程序退出 关闭链接
				return
			}
		case message, ok := <-c.Error: // 这个是有错误信息
			if !ok {
				// 发送关闭提示
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			res, _ := json.Marshal(map[string]interface{}{
				"_code":    100,
				"_message": message,
				"_data":    nil,
			})
			err := c.socket.WriteMessage(websocket.TextMessage, res)
			if err != nil {
				// 程序退出 关闭链接
				return
			}
		}
	}
}

// 客户端发送数据处理逻辑
func (c *Client) read() {
	defer func() {
		websocketManager.unregister <- c
		c.socket.Close()
	}()

	for {
		// 监听从 socket 获取数据
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			// 数据获取错误 退出登录
			fmt.Println("read 关闭")
			return
		}
		messageStruct := &protobuf.Message{}
		// proto解码
		proto.Unmarshal(message, messageStruct)
		// 找到对应的协议操作
		//fmt.Println(newGameMessage.Data)
		//buf := bytes.NewBuffer(newGameMessage.Data)
		//binary.Read(buf, binary.BigEndian, &newTestMessage)
		//newTestMessage = *(**protobuf.Test)(unsafe.Pointer(&newGameMessage.Data))
		//ptypes.UnmarshalAny(newGameMessage.Data, newTestMessage)
		var websocketReadMessage *ReadMessage
		//json.Unmarshal(message, &websocketReadMessage)
		if _, ok := ProtocolActions[messageStruct.Protocol]; !ok {
			c.Send <- []byte("请求错误")
		} else {
			//websocketReadMessage.Message = string(message[:])
			dataMessage := ProtocolStruct[messageStruct.Protocol]
			proto.Unmarshal(messageStruct.Data, dataMessage)
			handleFun, _ := ProtocolActions[messageStruct.Protocol]
			handleFun(&dataMessage, c, websocketReadMessage)
		}
		//激活start 程序 入广播管道
		//websocketManager.broadcast <- message
	}
}
