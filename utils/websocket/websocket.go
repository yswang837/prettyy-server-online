package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	WebSocketMap sync.Map
	OnlineCount  int32
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(toUser, message string) {
	if conn, ok := WebSocketMap.Load(toUser); ok {
		if err := conn.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("Error sending message to %s: %s\n", toUser, err)
		}
	}
}

func BroadcastMessage(message string) {
	WebSocketMap.Range(func(key, value interface{}) bool {
		go func(conn *websocket.Conn) {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Printf("Error broadcasting message: %s\n", err)
			}
		}(value.(*websocket.Conn))
		return true
	})
}

func BroadcastUserEvent(user string, event string) {
	currentUsers := GetCurrentUserList()
	message := fmt.Sprintf("%s%s聊天室，当前在线人数：%d。\n在线用户列表：%s\n", user, event, GetOnlineCount(), currentUsers)
	BroadcastMessage(message)
}

func CloseAndDeleteConnection(user string, conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println("Error closing connection: ", err)
	}
	WebSocketMap.Delete(user)
	reduceOnlineCount()
	BroadcastUserEvent(user, "离开了")
}

func AddOnlineCount() int32 {
	return atomic.AddInt32(&OnlineCount, 1)
}

func GetOnlineCount() int32 {
	return atomic.LoadInt32(&OnlineCount)
}

func reduceOnlineCount() int32 {
	return atomic.AddInt32(&OnlineCount, -1)
}

func GetCurrentUserList() string {
	users := make([]string, 0, 1000)
	WebSocketMap.Range(func(key, value interface{}) bool {
		user := key.(string)
		users = append(users, user)
		return true
	})
	userList, _ := json.Marshal(users)
	return string(userList)
}
