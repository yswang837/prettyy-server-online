package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// checkPasswordParams 面向接口
type checkPasswordParams struct {
	Uid string `json:"uid" form:"uid" binding:"required"`
}

// CheckPassword 检查密码是否为空，为空则在用户通过验证码登录后提示用户
// 4000020
// 2000020
func (s *Server) CheckPassword(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("check-password", "total")
	p := &checkPasswordParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("check-password", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000020, Message: "参数错误"})
		return
	}
	ctx.SetUid(p.Uid)
	u, err := user.GetUser(p.Uid)
	if err != nil {
		metrics.CommonCounter.Inc("check-password", "get-user-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000021, Message: "获取用户信息失败"})
		return
	}
	if u.Password == "" {
		metrics.CommonCounter.Inc("check-password", "empty-password")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000022, Message: "密码为空，请设置密码"})
		return
	}
	metrics.CommonCounter.Inc("check-password", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000020, Message: "有效的密码"})
	return
}
