package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"prettyy-server-online/utils/websocket"
)

func (s *Server) Lick(ctx *gin.Context) {
	websocket.Message <- "有人点赞了！"
	ctx.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}
