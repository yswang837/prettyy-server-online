package base

import (
	"net/http"
	xzfEmail "prettyy-server-online/custom-pkg/xzf-email"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// GetIdentifyCodeByEmail 通过邮件获取验证码，并将其写入到redis，5分钟过期
// 4000040
// 2000040
func (s *Server) GetIdentifyCodeByEmail(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("email-captcha", "total")
	email, _ := ctx.GetQuery("email")
	if email == "" {
		metrics.CommonCounter.Inc("email-captcha", "params-error")
		ctx.SetError("参数错误")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000040, Message: "参数错误"})
		return
	}
	ctx.SetEmail(email)
	//测试阶段不真正发邮件，该包已测试可用，可通过redis直接查看验证码
	//iCode := "667788"
	iCode, err := xzfEmail.SendEmail(email)
	if err != nil {
		metrics.CommonCounter.Inc("email-captcha", "send-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000041, Message: "邮件发送失败，请稍后重试或联系客服"})
		return
	}
	if err = user.SetExByEmail(email, iCode); err != nil {
		metrics.CommonCounter.Inc("email-captcha", "set-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000042, Message: "系统内部错误，请稍后重试或联系客服"})
		return
	}
	metrics.CommonCounter.Inc("email-captcha", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000040, Message: "设置邮箱验证码成功"})
	return
}
