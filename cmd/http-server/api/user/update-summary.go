package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updateNickNameParams 面向接口
type updateSummaryParams struct {
	Email   string `json:"email" form:"email" binding:"required"`
	Summary string `json:"summary" form:"summary" binding:"required"`
}

// UpdateSummary 更新用户的个人介绍
// 4000240
// 2000240
func (s *Server) UpdateSummary(ctx *gin.Context) {
	p := &updateSummaryParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000240, Message: "参数错误"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000241, Message: "获取用户信息失败"})
		return
	}
	if p.Summary == u.Summary {
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000242, Message: "个人简介未改变"})
		return
	}
	if err := user.UpdateSummary(p.Email, p.Summary); err != nil {
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000243, Message: "个人简介更新失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000240, Message: "个人简介更新成功"})
	return
}
