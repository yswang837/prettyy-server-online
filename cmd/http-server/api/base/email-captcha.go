package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	xzfEmail "prettyy-server-online/custom-pkg/xzf-email"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// GetIdentifyCodeByEmail 通过邮件获取验证码，并将其写入到redis，5分钟过期
// 4000040
// 2000040
func (s *Server) GetIdentifyCodeByEmail(ctx *gin.Context) {
	var err error
	email, _ := ctx.GetQuery("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000040, Message: "参数错误"})
		return
	}
	//测试阶段不真正发邮件，该包已测试可用，可通过redis直接查看验证码
	//iCode := "667788"
	iCode, err := xzfEmail.SendEmail(email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000041, Message: "邮件发送失败，请稍后重试或联系客服"})
		return
	}
	if err = user.SetExByEmail(email, iCode); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000042, Message: "设置邮箱验证码失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000040, Message: "设置邮箱验证码成功"})
	return
}
