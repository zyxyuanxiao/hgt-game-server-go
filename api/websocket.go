package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server/dto"
	"server/util"
	"time"
)

type WebsocketApi struct{}

type WebsocketClientManager struct {
	clients    map[int64]*WebsocketClient
	broadcast  chan []byte
	register   chan *WebsocketClient
	unregister chan *WebsocketClient
}

type WebsocketClient struct {
	UserDTO       dto.UserDTO
	socket        *websocket.Conn
	send          chan interface{}
	error         chan interface{}
	InsertTime    string
	GameGroupInfo struct {
		Status uint8
	}
	RoomId string
}

var websocketManager = WebsocketClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *WebsocketClient),
	unregister: make(chan *WebsocketClient),
	clients:    make(map[int64]*WebsocketClient),
}

type WebsocketReadMessage struct {
	Message string
	Action  string
	Args    interface{}
}

var WebsocketClientActions = make(map[string]func(*WebsocketClient, *WebsocketReadMessage))

// 加载映射方法
func LoadWebsocketActions() {
	WebsocketClientActions["createGroup"] = (*WebsocketClient).CreateGroup
	WebsocketClientActions["joinGroup"] = (*WebsocketClient).JoinGroup
	WebsocketClientActions["leaveGroup"] = (*WebsocketClient).LeaveGroup
}

func (websocketApi WebsocketApi) Websocket(c *gin.Context) {
	LoadWebsocketActions()
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
	client := &WebsocketClient{
		UserDTO:    userDto,
		socket:     conn,
		send:       make(chan interface{}),
		error:      make(chan interface{}),
		InsertTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	//把这个对象发送给 管道
	websocketManager.register <- client
	// 协程接收输出信息
	go client.write()
	go client.read()
}

func (manager *WebsocketClientManager) start() {
	for {
		select {
		case conn := <-manager.register: //新客户端加入
			manager.clients[conn.UserDTO.UserId] = conn
			fmt.Println("新用户加入:"+conn.UserDTO.Username, "加入时间："+conn.InsertTime)
			fmt.Println("当前总用户数量register：", len(manager.clients))
		case conn := <-manager.unregister: // 离开
			if _, ok := manager.clients[conn.UserDTO.UserId]; ok {
				close(conn.send)
				delete(manager.clients, conn.UserDTO.UserId)
				fmt.Println("用户离开：" + conn.UserDTO.Username)
				fmt.Println("当前总用户数量unregister：", len(manager.clients))
			}
		case message := <-manager.broadcast: //读到广播管道数据后的处理
			fmt.Println("当前总用户数量broadcast：", len(manager.clients))
			for _, conn := range manager.clients {
				select {
				case conn.send <- message: //调用发送给全体客户端
				default:
					// 重新上来之后挤掉了 @todo
					// 关闭连接
					close(conn.send)
					delete(manager.clients, conn.UserDTO.UserId)
				}
			}
		}
	}
}

// 广播数据 除了自己
func (manager *WebsocketClientManager) send(message []byte, ignore *WebsocketClient) {
	for _, conn := range manager.clients {
		if conn != ignore {
			conn.send <- message //发送的数据写入所有的 websocket 连接 管道
		}
	}
}

// 写入管道后激活这个进程
func (c *WebsocketClient) write() {
	defer func() {
		// 程序退出 关闭链接
		websocketManager.unregister <- c
		c.socket.Close()
	}()

	for {
		select {
		case data, ok := <-c.send: //这个管道有了数据 写这个消息出去
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
		case message, ok := <-c.error: // 这个是有错误信息
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
func (c *WebsocketClient) read() {
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
		var websocketReadMessage *WebsocketReadMessage
		json.Unmarshal(message, &websocketReadMessage)
		if _, ok := WebsocketClientActions[websocketReadMessage.Action]; !ok {
			c.send <- []byte("请求错误")
		} else {
			websocketReadMessage.Message = string(message[:])
			handleFun, _ := WebsocketClientActions[websocketReadMessage.Action]
			handleFun(c, websocketReadMessage)
		}
		//激活start 程序 入广播管道
		//websocketManager.broadcast <- message
	}
}
