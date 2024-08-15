package user

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updateGenderParams 面向接口
type updateGenderParams struct {
	Email  string `json:"email" form:"email" binding:"required"`
	Gender string `json:"gender" form:"gender" binding:"required"`
}

// UpdateGender 更新性别，默认只能更新一次
// 4000220
// 2000220
func (s *Server) UpdateGender(ctx *gin.Context) {
	p := &updateGenderParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000220, Message: "bind params err"})
		return
	}
	if err := user.UpdateGender(p.Email, p.Gender); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000221, Message: "update gender err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000220, Message: "update gender succ"})
	return
}
