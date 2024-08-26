package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000220, Message: "参数错误"})
		return
	}
	if err := user.UpdateGender(p.Email, p.Gender); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000221, Message: "更新性别失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000220, Message: "更新性别成功"})
	return
}
