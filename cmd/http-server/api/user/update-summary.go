package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// updateNickNameParams 面向接口
type updateSummaryParams struct {
	Uid     string `json:"uid" form:"uid" binding:"required"`
	Summary string `json:"summary" form:"summary" binding:"required"`
}

// UpdateSummary 更新用户的个人介绍
// 4000240
// 2000240
func (s *Server) UpdateSummary(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("update-summary", "total")
	p := &updateSummaryParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("update-summary", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000240, Message: "参数错误"})
		return
	}
	ctx.SetUid(p.Uid).SetSummary(p.Summary)
	u, err := user.GetUser(p.Uid)
	if err != nil {
		metrics.CommonCounter.Inc("update-summary", "get-user-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000241, Message: "获取用户信息失败"})
		return
	}
	if p.Summary == u.Summary {
		metrics.CommonCounter.Inc("update-summary", "same-summary")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000242, Message: "个人简介未改变"})
		return
	}
	if err = user.UpdateSummary(p.Uid, p.Summary); err != nil {
		metrics.CommonCounter.Inc("update-summary", "update-summary-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000243, Message: "个人简介更新失败"})
		return
	}
	metrics.CommonCounter.Inc("update-summary", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000240, Message: "个人简介更新成功"})
	return
}
