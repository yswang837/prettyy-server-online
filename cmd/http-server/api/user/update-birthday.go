package user

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updateBirthdayParams 面向接口
type updateBirthdayParams struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Birthday string `json:"birthday" form:"birthday" binding:"required"`
}

// UpdateBirthdayCity 更新用户的出生日期
// 4000280
// 2000280
func (s *Server) updateBirthdayParams(ctx *gin.Context) {
	p := &updateBirthdayParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000280, Message: "bind params err"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000281, Message: "get user err"})
		return
	}
	if p.Birthday == u.Birthday {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000282, Message: "birthday is same"})
		return
	}
	if err := user.UpdateBirthdayCity(p.Email, p.Birthday); err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000283, Message: "update birthday err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000280, Message: "update birthday succ"})
	return
}
