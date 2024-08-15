package user

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// updateNickNameParams 面向接口
type updateNickNameParams struct {
	Email    string `json:"email" form:"email" binding:"required"`
	NickName string `json:"nick_name" form:"nick_name" binding:"required"`
}

// UpdateNickName 更新用户名，默认用户名为邮箱前缀
// 4000200
// 2000200
func (s *Server) UpdateNickName(ctx *gin.Context) {
	p := &updateNickNameParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000200, Message: "bind params err"})
		return
	}
	u, err := user.GetUser(p.Email)
	if err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000201, Message: "get user err"})
		return
	}
	if p.NickName == u.NickName {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000202, Message: "nick name is same"})
		return
	}
	if err := user.UpdateNickName(p.Email, p.NickName); err != nil {
		ctx.JSON(200, ginConsulRegister.Response{Code: 4000203, Message: "update nick name err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000200, Message: "update nick name succ"})
	return
}
