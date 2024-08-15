package common

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	middle_ware "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
)

// Server 绑定所有通用的服务
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() (err error) {
	return nil
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middle_ware.JwtAuth())
	groupHandler.POST(conf.URLFileUpload, func(context *gin.Context) {
		s.FileUpload(context)
	})
}
