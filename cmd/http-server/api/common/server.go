package common

import (
	"github.com/gin-gonic/gin"
	"os"
	"prettyy-server-online/cmd/http-server/conf"
	middle_ware "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
	"sync"
)

var (
	uploadDir = ""
)

// Server 绑定所有通用的服务
type Server struct {
	sseChannelMap *sync.Map
}

func NewServer() *Server {
	return &Server{}
}

// Init 初始化服务，该目录是nginx托管的前端静态文件的目录，后端upload上来的文件直接放到这里就可以直接访问了
func (s *Server) Init() (err error) {
	// 初始化sse通道
	s.sseChannelMap = &sync.Map{}
	// 初始化上传文件的目录
	return os.MkdirAll(uploadDir, os.ModePerm)
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 本地测试图片上传的相关的静态目录，仅仅用于本地测试，测试时需关闭auth中间件
	r.Static("/uploads", "/Users/yuanshun/workspace/my/prettyy-server-online/uploads")
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middle_ware.JwtAuth())
	groupHandler.POST(conf.URLFileUpload, func(context *gin.Context) {
		s.FileUpload(context)
	})
	groupHandler.GET(conf.URLSseConnection, func(context *gin.Context) {
		s.SseConnection(context)
	})
	groupHandler.GET(conf.URLSseSend, func(context *gin.Context) {
		s.SseSend(context)
	})
}

func init() {
	if os.Getenv("idc") == "dev" {
		uploadDir = "./uploads"
	} else {
		// 服务器nginx托管的静态资源目录
		uploadDir = "/root/prettyy-web-online/dist/uploads"
	}
}
