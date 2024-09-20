package common

import (
	"github.com/gin-gonic/gin"
	"os"
	"prettyy-server-online/cmd/http-server/conf"
	middle_ware "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
)

const (
	uploadDir = "/root/prettyy-web-online/dist/uploads"
)

// Server 绑定所有通用的服务
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

// Init 初始化服务，该目录是nginx托管的前端静态文件的目录，后端upload上来的文件直接放到这里就可以直接访问了
func (s *Server) Init() (err error) {
	if os.Getenv("idc") == "dev" {
		return os.MkdirAll("./uploads", os.ModePerm)
	}
	return os.MkdirAll(uploadDir, os.ModePerm)
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middle_ware.JwtAuth())
	groupHandler.POST(conf.URLFileUpload, func(context *gin.Context) {
		s.FileUpload(context)
	})
}
