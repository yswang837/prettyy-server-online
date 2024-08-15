package base

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
)

// GetIdentifyCodeByEmail 通过邮件获取验证码，并将其写入到redis，5分钟过期
// 4000040
// 2000040
func (s *Server) GetIdentifyCodeByEmail(ctx *gin.Context) {
	var err error
	email, _ := ctx.GetQuery("email")
	//测试阶段不真正发邮件，该包已测试可用，可通过redis直接查看验证码
	//iCode, err := xzfEmail.SendEmail(email)
	//if err != nil {
	//	ctx.JSON(200, ginConsulRegister.Response{Code: "4000040", Msg: "send email err"})
	//	return
	//}
	iCode := "667788"
	if err = user.SetExByEmail(email, iCode); err != nil {
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000041, Message: "set identify code err"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000040, Message: "set identify code success"})
	return
}
