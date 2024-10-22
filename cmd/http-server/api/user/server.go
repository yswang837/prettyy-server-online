package user

import (
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	middleWare "prettyy-server-online/custom-pkg/xzf-gin-consul/middle-ware"
	"prettyy-server-online/custom-pkg/xzf-gin-consul/register"
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
	groupHandler := r.Group("").Use(middleWare.JwtAuth())
	groupHandler.GET(conf.URLLoginOut, func(context *gin.Context) {
		s.LoginOut(register.NewContext(context))
	})
	groupHandler.GET(conf.URLCheckPassword, func(context *gin.Context) {
		s.CheckPassword(register.NewContext(context))
	})
	groupHandler.POST(conf.URLUpdateNickName, func(context *gin.Context) {
		s.UpdateNickName(register.NewContext(context))
	})
	groupHandler.POST(conf.URLUpdateGender, func(context *gin.Context) {
		s.UpdateGender(register.NewContext(context))
	})
	groupHandler.POST(conf.URLUpdateSummary, func(context *gin.Context) {
		s.UpdateSummary(register.NewContext(context))
	})
	groupHandler.POST(conf.URLUpdateProvinceCity, func(context *gin.Context) {
		s.UpdateProvinceCity(register.NewContext(context))
	})
	groupHandler.POST(conf.URLUpdateBirthday, func(context *gin.Context) {
		s.UpdateBirthday(register.NewContext(context))
	})
	groupHandler.POST(conf.URLUpdatePassword, func(context *gin.Context) {
		s.UpdatePassword(register.NewContext(context))
	})

}
