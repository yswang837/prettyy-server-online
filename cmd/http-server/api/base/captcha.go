package base

import (
	"github.com/mojocn/base64Captcha"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/utils/metrics"
)

type captchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
}

var store = base64Captcha.DefaultMemStore

// GetIdentifyCode 账密登录时，获取验证码
// 4000100
// 2000100
func (s *Server) GetIdentifyCode(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("captcha", "total")
	driver := base64Captcha.NewDriverDigit(80, 180, 4, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		metrics.CommonCounter.Inc("captcha", "err")
		ctx.SetError("生成验证码失败")
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000100, Message: "生成验证码失败，请稍后重试或联系客服"})
		return
	}
	resp := &captchaResponse{CaptchaId: id, PicPath: b64s, CaptchaLength: 4}
	metrics.CommonCounter.Inc("captcha", "succ")
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000100, Message: "生成验证码成功", Result: resp})
	return
}
