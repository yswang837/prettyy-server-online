package user

import (
	"github.com/gin-gonic/gin"
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
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000240, Message: "bind params err"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000241, Message: "get user err"})
		return
	}
	if p.Summary == u.Summary {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000242, Message: "summary is same"})
		return
	}
	if err := user.UpdateSummary(p.Email, p.Summary); err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000243, Message: "update summary err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000240, Message: "update summary succ"})
	return
}
