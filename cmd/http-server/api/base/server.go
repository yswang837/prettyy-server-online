package base

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	middleWare "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
)

// Server 绑定所有用户相关的服务
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() (err error) {
	return nil
}

func (s *Server) SetRoute(r *gin.Engine) {
	r.POST(conf.URLRegisterLogin, func(context *gin.Context) {
		s.LoginRegister(context)
	})
	// 通过邮件发送验证码，账密登录获取验证码，都需要走频次限制的中间件
	groupHandler := r.Group("").Use(middleWare.Restrict())
	groupHandler.GET(conf.URLIdentifyCodeByEmail, func(context *gin.Context) {
		s.GetIdentifyCodeByEmail(context)
	})
	groupHandler.GET(conf.URLIdentifyCode, func(context *gin.Context) {
		s.GetIdentifyCode(context)
	})
}
