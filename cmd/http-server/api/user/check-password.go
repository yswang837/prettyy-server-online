package user

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// checkPasswordParams 面向接口
type checkPasswordParams struct {
	Email string `json:"email" form:"email" binding:"required"`
}

// CheckPassword 检查密码是否为空，为空则在用户通过验证码登录后提示用户
// 4000020
// 2000020
func (s *Server) CheckPassword(ctx *gin.Context) {
	p := &checkPasswordParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000020, Message: "bind params err"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000021, Message: "get user err"})
		return
	}
	if u.Password == "" {
		ctx.JSON(200, ginConsulRegister.Response{Code: 2000020, Message: "empty password, please set it"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000021, Message: "valid password"})
	return
}
