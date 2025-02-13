package service

import (
	"chat/conf"
	"chat/pkg/e"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("---监听管道通信---")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("---有新连接---", conn.ID)
			fmt.Println("Manager.Register", conn)
			//把连接放到用户管理上
			Manager.Clients[conn.ID] = conn
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已经连接到服务器了",
			}
			msg, err := json.Marshal(replyMsg) //序列化
			if err != nil {
				log.Printf("序列化消息失败: %v", err)
				continue
			}
			conn.mu.Lock() // 加锁
			err = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			conn.mu.Unlock() // 解锁
			if err != nil {
				log.Printf("发送消息失败: %v", err)
			}
		case conn := <-Manager.Unregister:
			fmt.Printf("---连接失败---", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接已经断开",
				}
				msg, err := json.Marshal(replyMsg) //序列化
				if err != nil {
					log.Printf("序列化消息失败: %v", err)
					continue
				}
				conn.mu.Lock() // 加锁
				err = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				conn.mu.Unlock() // 解锁
				if err != nil {
					log.Printf("发送消息失败: %v", err)
				}
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast: //1->2
			message := broadcast.Message
			sendId := broadcast.Client.SendID //2->1
			flag := false                     //默认对方不在线
			for id, conn := range Manager.Clients {
				if id == sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID //1->2
			if flag {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线",
				}
				msg, err := json.Marshal(replyMsg) //序列化
				if err != nil {
					log.Printf("序列化消息失败: %v", err)
					continue
				}
				broadcast.Client.mu.Lock() // 加锁
				err = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				broadcast.Client.mu.Unlock() // 解锁
				if err != nil {
					log.Printf("发送消息失败: %v", err)
				}
				log.Printf("准备插入消息到 MongoDB，数据库名称: %s", conf.MongoDBName)
				err = InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) //已经读了
				if err != nil {
					fmt.Println("插入消息失败", err)
				}
			} else {
				fmt.Println("对方不在线")
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线",
				}
				msg, err := json.Marshal(replyMsg) //序列化
				if err != nil {
					log.Printf("序列化消息失败: %v", err)
					continue
				}
				broadcast.Client.mu.Lock() // 加锁
				err = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				broadcast.Client.mu.Unlock() // 解锁
				if err != nil {
					log.Printf("发送消息失败: %v", err)
				}
				log.Printf("准备插入消息到 MongoDB，数据库名称: %s", conf.MongoDBName)
				err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					fmt.Println("插入消息失败", err)
				}
			}
		}
	}
}
