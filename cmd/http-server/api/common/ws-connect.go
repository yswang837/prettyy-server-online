package common

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/utils/websocket"
)

func (s *Server) WSConnect(ctx *gin.Context) {
	ws, err := websocket.UpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	s.mu.Lock()
	websocket.Clients[ws] = true
	s.mu.Unlock()
	for {
		_, _, err = ws.ReadMessage()
		if err != nil {
			s.mu.Lock()
			delete(websocket.Clients, ws)
			s.mu.Unlock()
			break
		}
	}
}
