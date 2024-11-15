package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// updatePasswordParams 面向接口
type updatePasswordParams struct {
	Uid      string `json:"uid" form:"uid" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// UpdatePassword 设置密码，设置后即可账密登录
// 4000300
// 2000300
func (s *Server) UpdatePassword(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("update-password", "total")
	p := &updatePasswordParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("update-password", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000300, Message: "参数错误"})
		return
	}
	ctx.SetUid(p.Uid).SetPassword(p.Password)
	// 密码长度控制在6~20位
	if len(p.Password) < 6 || len(p.Password) > 20 {
		metrics.CommonCounter.Inc("update-password", "invalid-password")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000301, Message: "密码长度必须在6~20个字符"})
		return
	}
	u, err := user.GetUser(p.Uid)
	if err != nil {
		metrics.CommonCounter.Inc("update-password", "get-user-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000302, Message: "获取用户信息失败"})
		return
	}
	if u.Password != "" && u.Password == p.Password {
		// 库中的密码不为空，且密码相同
		metrics.CommonCounter.Inc("update-password", "same-password")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000303, Message: "密码未改变"})
		return
	}
	// 要么是库中密码为空，要么是密码不同，均可以直接更新密码，binding required 已经保证了密码不为空
	if err = user.UpdatePassword(p.Uid, p.Password); err != nil {
		metrics.CommonCounter.Inc("update-password", "update-password-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000304, Message: "密码更新失败"})
		return
	}
	metrics.CommonCounter.Inc("update-password", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000300, Message: "密码更新成功"})
	return
}
