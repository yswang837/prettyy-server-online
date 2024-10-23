package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// updateBirthdayParams 面向接口
type updateBirthdayParams struct {
	Uid      string `json:"uid" form:"uid" binding:"required"`
	Birthday string `json:"birthday" form:"birthday" binding:"required"`
}

// UpdateBirthday 更新用户的出生日期
// 4000280
// 2000280
func (s *Server) UpdateBirthday(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("update-birthday", "total")
	p := &updateBirthdayParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("update-birthday", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000280, Message: "参数错误"})
		return
	}
	ctx.SetUid(p.Uid).SetBirthday(p.Birthday)
	u, err := user.GetUser(p.Uid)
	if err != nil {
		metrics.CommonCounter.Inc("update-birthday", "get-user-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000281, Message: "获取用户信息失败"})
		return
	}
	if p.Birthday == u.Birthday {
		metrics.CommonCounter.Inc("update-birthday", "same-birthday")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000282, Message: "生日未改变"})
		return
	}
	if err = user.UpdateBirthdayCity(p.Uid, p.Birthday); err != nil {
		metrics.CommonCounter.Inc("update-birthday", "update-birthday-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000283, Message: "更新生日失败"})
		return
	}
	metrics.CommonCounter.Inc("update-birthday", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000280, Message: "更新生日成功"})
	return
}
