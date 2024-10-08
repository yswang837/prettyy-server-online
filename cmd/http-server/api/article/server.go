package article

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"prettyy-server-online/cmd/http-server/conf"
	middleWare "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
	xzfSnowflake "prettyy-server-online/custom-pkg/xzf-snowflake"
	"prettyy-server-online/utils/http"
)

// Server 绑定所有文章相关的服务
type Server struct {
	client *resty.Client
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Init() (err error) {
	s.client = http.NewClient()
	return xzfSnowflake.Init("2024-03-09", "1")
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 不需要token认证
	r.GET(conf.URLGetArticleDetail, func(context *gin.Context) {
		s.ArticleDetail(context)
	})
	r.GET(conf.URLGetArticleList, func(context *gin.Context) {
		s.ArticleList(context)
	})
	r.GET(conf.URLGetUserInfoByAid, func(context *gin.Context) {
		s.GetUserInfoByAid(context)
	})
	r.POST(conf.URLExtractSummary, func(context *gin.Context) {
		s.ExtractSummary(context)
	})
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middleWare.JwtAuth())
	groupHandler.POST(conf.URLPublishArticle, func(context *gin.Context) {
		s.PublishArticle(context)
	})
	groupHandler.POST(conf.URLDelArticle, func(context *gin.Context) {
		s.DelArticle(context)
	})
}
