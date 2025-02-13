package service

import (
	"chat/pkg/e"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("---监听管道通信---")
		select {
		case conn := <-Manager.Register:
			fmt.Println("---有新连接---", conn.ID)
			//把连接放到用户管理上
			Manager.Clients[conn.ID] = conn
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已经连接到服务器了",
			}
			msg, _ := json.Marshal(replyMsg) //序列化
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
