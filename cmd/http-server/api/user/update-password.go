package user

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updatePasswordParams 面向接口
type updatePasswordParams struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// UpdatePassword 设置密码，设置后即可账密登录
// 4000300
// 2000300
func (s *Server) UpdatePassword(ctx *gin.Context) {
	p := &updatePasswordParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000300, Message: "bind params err"})
		return
	}
	// 密码长度控制在6~20位
	if len(p.Password) < 6 || len(p.Password) > 20 {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000301, Message: "password length must be 6~20"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000302, Message: "get user err"})
		return
	}
	if u.Password != "" && u.Password == p.Password {
		// 库中的密码不为空，且密码相同
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000303, Message: "password is same"})
		return
	}
	// 要么是库中密码为空，要么是密码不同，均可以直接更新密码，binding required 已经保证了密码不为空
	if err = user.UpdatePassword(p.Email, p.Password); err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000304, Message: "update password err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000300, Message: "update password success"})
	return
}
