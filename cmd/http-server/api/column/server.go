package column

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
)

// Server 绑定所有专栏相关的服务
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() (err error) {
	return
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 不需要token认证，文章详情页获取专栏列表不需要认证，专栏管理获取专栏列表需要认证
	r.GET(conf.URLGetColumnList, func(context *gin.Context) {
		s.ColumnList(context)
	})
}
