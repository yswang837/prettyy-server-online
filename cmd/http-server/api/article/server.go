package article

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"prettyy-server-online/cmd/http-server/conf"
	middleWare "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
	"prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	xzfSnowflake "prettyy-server-online/custom-pkg/xzf-snowflake"
	"prettyy-server-online/utils/http"
	"time"
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
	s.client.SetTimeout(time.Second * 20) // 调用通义千问模型提取摘要，比较慢，设置20秒超时
	return xzfSnowflake.Init("2024-03-09", "1")
}

func (s *Server) SetRoute(r *gin.Engine) {
	// 不需要token认证
	r.GET(conf.URLGetArticleDetail, func(context *gin.Context) {
		s.ArticleDetail(register.NewContext(context))
	})
	r.GET(conf.URLGetArticleList, func(context *gin.Context) {
		s.ArticleList(register.NewContext(context))
	})
	r.GET(conf.URLGetUserInfoByAid, func(context *gin.Context) {
		s.GetUserInfoByAid(register.NewContext(context))
	})
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middleWare.JwtAuth())
	groupHandler.POST(conf.URLPublishArticle, func(context *gin.Context) {
		s.PublishArticle(register.NewContext(context))
	})
	groupHandler.POST(conf.URLDelArticle, func(context *gin.Context) {
		s.DelArticle(register.NewContext(context))
	})
	groupHandler.POST(conf.URLLikeCollectArticle, func(context *gin.Context) {
		s.ClickLikeCollect(register.NewContext(context))
	})
	groupHandler.POST(conf.URLExtractSummary, func(context *gin.Context) {
		s.ExtractSummary(register.NewContext(context))
	})
}
