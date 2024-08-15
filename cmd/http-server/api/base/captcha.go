package base

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
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
func (s *Server) GetIdentifyCode(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000100, Message: "生成验证码失败"})
		return
	}
	resp := &captchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: 6,
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000100, Message: "生成验证码成功", Result: resp})
	return
}
