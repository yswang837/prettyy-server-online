package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// updateGenderParams 面向接口
type updateGenderParams struct {
	Uid    string `json:"uid" form:"uid" binding:"required"`
	Gender string `json:"gender" form:"gender" binding:"required"`
}

// UpdateGender 更新性别，默认只能更新一次
// 4000220
// 2000220
func (s *Server) UpdateGender(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("update-gender", "total")
	p := &updateGenderParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("update-gender", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000220, Message: "参数错误"})
		return
	}
	ctx.SetUid(p.Uid).SetGender(p.Gender)
	if err := user.UpdateGender(p.Uid, p.Gender); err != nil {
		metrics.CommonCounter.Inc("update-gender", "update-gender-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000221, Message: "更新性别失败"})
		return
	}
	metrics.CommonCounter.Inc("update-gender", "succ")
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000220, Message: "更新性别成功"})
	return
}
