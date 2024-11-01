package websocket

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	middle_ware "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
	"prettyy-server-online/custom-pkg/xzf-gin-consul/register"
)

// Server 绑定所有通用的服务
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// Init 初始化服务，该目录是nginx托管的前端静态文件的目录，后端upload上来的文件直接放到这里就可以直接访问了
func (s *Server) Init() (err error) {
	return nil
}

func (s *Server) SetRoute(r *gin.Engine) {
	groupHandler := r.Group("").Use(middle_ware.JwtAuth())
	groupHandler.GET(conf.URLChat, func(context *gin.Context) {
		s.ChatHandler(register.NewContext(context))
	})
}
