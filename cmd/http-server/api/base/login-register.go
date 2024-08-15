package base

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	user2 "prettyy-server-online/data/user"
	user3 "prettyy-server-online/services/user"
	"prettyy-server-online/utils/tool"
)

// loginRegisterParams 面向接口
type loginRegisterParams struct {
	Email        string `json:"email" form:"email" binding:"required"`   // 邮箱：目前支持这两种登录方式，后期新增一个微信扫码登录的方式
	Method       string `json:"method" form:"method" binding:"required"` // 登录或注册的方式：1验证码登录 2密码登录，在前端第一个选项卡值为1
	Password     string `json:"password" form:"password"`                // 密码：验证码登录方式不需要密码
	IdentifyCode string `json:"identify_code" form:"identify_code"`      // 验证码：密码登录方式不需要验证码
	IdentifyID   string `json:"identify_id" form:"identify_id"`          // 验证码ID：密码登录方式不需要验证码ID
}

// LoginRegister 登录或注册接口，提示用户：如果未注册，那么登录时将自动注册
// 4000001~4000009
// 2000001~2000003
// todo 计数和详细日志
func (s *Server) LoginRegister(ctx *gin.Context) {
	p := &loginRegisterParams{}
	var err error
	if err = ctx.Bind(p); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000001, Message: "参数绑定错误"})
		return
	}
	switch p.Method {
	case "1":
		// 验证码登录/注册
		if p.IdentifyCode == "" {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000002, Message: "验证码为空"})
			return
		}
		if p.IdentifyCode != user3.GetIdentifyCodeFromCache(p.Email) {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000003, Message: "验证码错误"})
			return
		}
	case "2":
		// 密码登录/注册
		if p.Password == "" {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000004, Message: "密码为空"})
			return
		}
		if p.IdentifyCode == "" {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000002, Message: "验证码为空"})
			return
		}
	default:
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000005, Message: "不支持的登录/注册方式"})
		return
	}
	// 执行到这，验证码通过 或者 密码不为空
	// 检查用户是否已经注册

	token, err := tool.GenerateToken()
	if err != nil || token == "" {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000010, Message: "生成token失败"})
		return
	}
	u, _ := user3.GetUser(p.Email)
	if u == nil {
		// 未注册
		user := &user2.User{Email: p.Email, Password: p.Password}
		if err = user3.Add(user); err != nil {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000006, Message: "注册失败"})
			return
		}
		if err = user3.UpdateLoginTime(p.Email); err != nil {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000007, Message: "更新登录时间失败"})
			return
		}
		user, err := user3.GetUser(p.Email)
		if err != nil {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000011, Message: "注册成功，但获取用户信息失败"})
		}
		result := map[string]interface{}{"token": token, "user": user}
		ctx.JSON(200, ginConsulRegister.Response{Code: 2000001, Message: "注册成功", Result: result})
		return
	} else {
		// 已注册，走登录逻辑
		user, err := user3.GetUser(p.Email)
		if err != nil {
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000011, Message: "注册成功，但获取用户信息失败"})
		}
		result := map[string]interface{}{"token": token, "user": user}
		switch p.Method {
		case "1":
			// 走到这里，验证码和密码都匹配了，登录成功更新登录时间
			if err = user3.UpdateLoginTime(p.Email); err != nil {
				ctx.JSON(400, ginConsulRegister.Response{Code: 4000007, Message: "更新登录时间失败"})
				return
			}
			ctx.JSON(200, ginConsulRegister.Response{Code: 2000002, Message: "验证码登录成功", Result: result})
			return
		case "2":
			if !store.Verify(p.IdentifyID, p.IdentifyCode, true) {
				// 验证码错误
				ctx.JSON(400, ginConsulRegister.Response{Code: 40000010, Message: "验证码错误"})
				return
			}
			if u.Password == tool.ToMd5(p.Password) {
				// 登录成功更新登录时间
				if err = user3.UpdateLoginTime(p.Email); err != nil {
					ctx.JSON(400, ginConsulRegister.Response{Code: 4000007, Message: "更新登录时间失败"})
					return
				}
				ctx.JSON(200, ginConsulRegister.Response{Code: 2000003, Message: "账密登录成功", Result: result})
				return
			} else {
				// 该情况可能为：
				if u.Password == "" {
					// 1、用户通过验证码注册，从而未设置密码(数据库中密码为空)，而登录的时候走了密码登录
					ctx.JSON(400, ginConsulRegister.Response{Code: 4000008, Message: "您未设置密码，请使用免密登录后设置密码"})
					return
				} else {
					// 2、用户输入的密码确实有误
					ctx.JSON(400, ginConsulRegister.Response{Code: 4000009, Message: "邮箱或者密码错误"})
					return
				}
			}
		default:
			ctx.JSON(400, ginConsulRegister.Response{Code: 4000005, Message: "不支持的登录/注册方式"})
			return
		}
	}
}
