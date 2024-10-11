package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
)

// SseConnection 建立sse连接 Server Sent Event
// 4000420
// 2000420

type sseConnectionParams struct {
	UserEmail string `json:"user_email" form:"user_email"`
	TraceId   string `json:"trace_id" form:"trace_id"`
}

func (s *Server) SseConnection(ctx *gin.Context) {
	params := &sseConnectionParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000420, Message: "参数错误"})
		return
	}
	newChannel := make(chan string)
	mKey := params.UserEmail + params.TraceId
	s.sseChannelMap.Store(mKey, newChannel)
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	// 监听客户端通道是否被关闭
	closeNotify := ctx.Request.Context().Done()
	go func() {
		<-closeNotify
		s.sseChannelMap.Delete(mKey)
		return
	}()

	// 获取http写入器并断言为flusher，让其将缓冲器的数据立即写入
	w := ctx.Writer
	flusher, _ := w.(http.Flusher)
	curChan, _ := s.sseChannelMap.Load(mKey)
	for msg := range curChan.(chan string) {
		fmt.Fprintf(w, "data:%s\n\n", msg)
		flusher.Flush()
	}
	return
}
