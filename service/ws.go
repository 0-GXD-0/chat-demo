package service

import (
	"chat/cache"
	"chat/pkg/e"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const month = 60 * 60 * 24 * 30 //一个月30天

type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
	mu     sync.Mutex //添加互斥锁
}

type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) //升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//创建用户实例
	client := &Client{
		ID:     CreateID(uid, toUid),
		SendID: CreateID(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	//用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)

		err := c.Socket.ReadJSON(sendMsg)
		if err != nil {
			fmt.Println("数据格式不正确", err)
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		if sendMsg.Type == 1 {
			//发送消息
			r1, _ := cache.RedisClient.Get(c.ID).Result()
			r2, _ := cache.RedisClient.Get(c.SendID).Result()
			if r1 > "3" && r2 == "" {
				//1给2发消息，但是2没有回，或没有看到，就停止1发送
				replyMsg := ReplyMsg{
					Code:    e.WebsocketLimit,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replyMsg) //序列化
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
		} else {
			cache.RedisClient.Incr(c.ID)
			_, _ = cache.RedisClient.Expire(c.ID, month).Result()
		}
		//传到广播
		Manager.Broadcast <- &Broadcast{
			Client:  c,
			Message: []byte(sendMsg.Content),
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", message),
			}
			msg, _ := json.Marshal(replyMsg) //序列化
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
