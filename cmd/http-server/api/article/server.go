package article

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	middle_ware "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
)

// Server 绑定所有文章相关的服务
type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() (err error) {
	return nil
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 不需要token认证
	r.GET(conf.URLGetArticleDetail, func(context *gin.Context) {
		s.ArticleDetail(context)
	})
	r.GET(conf.URLGetArticleList, func(context *gin.Context) {
		s.ArticleList(context)
	})
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middle_ware.JwtAuth())
	groupHandler.POST(conf.URLPublishArticle, func(context *gin.Context) {
		s.PublishArticle(context)
	})
	groupHandler.POST(conf.URLDelArticle, func(context *gin.Context) {
		s.DelArticle(context)
	})
}
