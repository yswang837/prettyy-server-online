package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"strings"
	"time"
)

// SseSend 服务端向客户端推送sse消息 Server Sent Event
// 4000440
// 2000440

type sseSendParams struct {
	UserEmail   string `json:"user_email" form:"user_email"`
	MessageBody string `json:"message_body" form:"message_body"`
	ActionType  string `json:"action_type" form:"action_type"`
}

type sseSendResp struct {
	MessageBody string    `json:"message_body"`
	UserEmail   string    `json:"user_email"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatTime   time.Time `json:"creat_time"`
}

func (s *Server) SseSend(ctx *gin.Context) {
	params := &sseSendParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000440, Message: "参数错误"})
		return
	}
	log.Print("Send notification to user = " + params.UserEmail)
	var msg = sseSendResp{
		MessageBody: params.MessageBody,
		UserEmail:   params.UserEmail,
		Type:        params.ActionType,
		Status:      "UNREAD",
		CreatTime:   time.Now(),
	}
	msgBytes, _ := json.Marshal(msg)
	s.sseChannelMap.Range(func(key, value any) bool {
		k := key.(string)
		if strings.Contains(k, params.UserEmail) {
			channel := value.(chan string)
			channel <- string(msgBytes)
		}
		return true
	})
	return
}
