package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/utils/websocket"
)

type Message struct {
	Content string `json:"content"`
	ToUser  string `json:"toUser"`
}

func (s *Server) ChatHandler(c *ginConsulRegister.Context) {
	user := c.Query("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名参数不能为空"})
		return
	}
	ws, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket 升级失败：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket 升级失败"})
		return
	}
	defer websocket.CloseAndDeleteConnection(user, ws)

	websocket.WebSocketMap.Store(user, ws)
	websocket.AddOnlineCount()
	websocket.BroadcastUserEvent(user, "加入了")
	log.Printf("用户连接：%s，当前在线人数：%d\n", user, websocket.GetOnlineCount())

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			break
		}
		var msg Message
		if err := json.Unmarshal(p, &msg); err != nil {
			websocket.SendMessage(user, "消息格式错误："+err.Error())
			continue
		}

		if msg.ToUser == "" || msg.Content == "" {
			websocket.SendMessage(user, "接收人或消息内容不能为空")
			continue
		} else if _, ok := websocket.WebSocketMap.Load(msg.ToUser); !ok {
			websocket.SendMessage(user, "接收人不在线")
			continue
		}

		msgFrom := fmt.Sprintf("来自%s的消息：%s", user, msg.Content)
		websocket.SendMessage(msg.ToUser, msgFrom)
	}
}
