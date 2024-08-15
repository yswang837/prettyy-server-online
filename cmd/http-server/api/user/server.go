package user

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	middle_ware "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
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
	// 需要token认证的路由组
	groupHandler := r.Group("").Use(middle_ware.JwtAuth())
	groupHandler.GET(conf.URLLoginOut, func(context *gin.Context) {
		s.LoginOut(context)
	})
	groupHandler.GET(conf.URLCheckPassword, func(context *gin.Context) {
		s.CheckPassword(context)
	})
	groupHandler.POST(conf.URLUpdateNickName, func(context *gin.Context) {
		s.UpdateNickName(context)
	})
	groupHandler.POST(conf.URLUpdateGender, func(context *gin.Context) {
		s.UpdateGender(context)
	})
	groupHandler.POST(conf.URLUpdateSummary, func(context *gin.Context) {
		s.UpdateSummary(context)
	})
	groupHandler.POST(conf.URLUpdateProvinceCity, func(context *gin.Context) {
		s.UpdateProvinceCity(context)
	})
	groupHandler.POST(conf.URLUpdateBirthday, func(context *gin.Context) {
		s.updateBirthdayParams(context)
	})
}
