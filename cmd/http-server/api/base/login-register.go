package base

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	invertedIndex "prettyy-server-online/data/inverted-index"
	user2 "prettyy-server-online/data/user"
	invertedIndex2 "prettyy-server-online/services/inverted-index"
	user3 "prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
	"prettyy-server-online/utils/tool"
	"strconv"
)

// loginRegisterParams 面向接口
type loginRegisterParams struct {
	Email        string `form:"email" binding:"required"`  // 邮箱：目前支持这两种登录方式，后期新增一个微信扫码登录的方式
	Method       string `form:"method" binding:"required"` // 登录或注册的方式：1验证码登录 2密码登录，在前端第一个选项卡值为1
	IdentifyCode string `form:"identify_code"`             // 验证码：账密登录时，验证码是前端填的那个图片验证码的数据，免密登录时，它是邮箱验证码的值
	// 免密登录仅需上面三个值，账密登录需要所有
	Password   string `form:"password"`    // 密码：仅账密登录方式需要密码
	IdentifyID string `form:"identify_id"` // 验证码ID：仅账密登录时需要验证码ID
}

// LoginRegister 登录或注册接口，提示用户：如果未注册，那么登录时将自动注册。接口返回的登录时间，是上一次登录的时间
// 4000001
// 2000001
func (s *Server) LoginRegister(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("login-register", "total")
	p := &loginRegisterParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("login-register", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000001, Message: "参数错误"})
		return
	}
	ctx.SetEmail(p.Email).SetMethod(p.Method).SetPassword(tool.ToMd5(p.Password))
	switch p.Method {
	case "1":
		metrics.CommonCounter.Inc("login-register", "captcha-1")
		// 验证码登录/注册
		if p.IdentifyCode == "" {
			metrics.CommonCounter.Inc("login-register", "empty-captcha-1")
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000002, Message: "验证码为空"}) // 免密方式，验证码为空
			return
		}
		if p.IdentifyCode != user3.GetIdentifyCodeFromCache(p.Email) {
			metrics.CommonCounter.Inc("login-register", "error-captcha-1")
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000003, Message: "验证码错误"}) //免密方式，验证码错误
			return
		}
	case "2":
		metrics.CommonCounter.Inc("login-register", "captcha-2")
		// 密码登录/注册
		if p.Password == "" {
			metrics.CommonCounter.Inc("login-register", "empty-password-2")
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000004, Message: "密码为空"}) //账密方式，密码为空
			return
		}
		if p.IdentifyID == "" || p.IdentifyCode == "" {
			metrics.CommonCounter.Inc("login-register", "empty-captcha-2")
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000005, Message: "验证码为空"}) //账密方式，验证码为空
			return
		}
		if !store.Verify(p.IdentifyID, p.IdentifyCode, true) {
			// 验证码错误，防爆次数为1，也就是填错了就清空当前的identify_id
			metrics.CommonCounter.Inc("login-register", "error-captcha-2")
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000006, Message: "验证码错误"}) //账密方式，验证码错误
			return
		}
	default:
		metrics.CommonCounter.Inc("login-register", "unsupported-method")
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000007, Message: "不支持的登录/注册方式"})
		return
	}
	// 执行到这，两种验证码都通过 或者 账密登录时密码不为空
	token := ""
	// 检查用户是否已经注册，通过Email检查反向索引库，如果存在，则已注册，否则未注册
	i, _ := invertedIndex2.Get(invertedIndex.TypEmailUid, p.Email)
	if len(i) == 0 {
		metrics.CommonCounter.Inc("login-register", "register")
		// 未注册，走注册逻辑
		user := &user2.User{Email: p.Email, Password: p.Password}
		userObj, err := user3.Add(user)
		if err == nil {
			// 添加反向索引
			invertedObj := &invertedIndex.InvertedIndex{AttrValue: p.Email, Typ: invertedIndex.TypEmailUid, Idx: strconv.FormatInt(userObj.Uid, 10)}
			if err = invertedIndex2.Add(invertedObj); err != nil {
				metrics.CommonCounter.Inc("login-register", "register-failed")
				ctx.SetError(err.Error())
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000009, Message: "注册失败"})
				return
			}
			// 生成token，无论注册还是登录均带上token返回
			token, err = tool.GenerateToken(userObj.Uid)
			if err != nil || token == "" {
				metrics.CommonCounter.Inc("login-register", "generate-token")
				if err != nil {
					ctx.SetError(err.Error())
				}
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000008, Message: "生成token失败"})
				return
			}
		} else {
			metrics.CommonCounter.Inc("login-register", "register-failed")
			ctx.SetError(err.Error())
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000009, Message: "注册失败"})
			return
		}
		m := user3.UserToMap(userObj)
		delete(m, "uid")
		result := map[string]interface{}{"token": token, "user": m}
		metrics.CommonCounter.Inc("login-register", "register-succ")
		ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000001, Message: "注册成功", Result: result})
		return
	} else {
		metrics.CommonCounter.Inc("login-register", "login")
		// 已注册，走登录逻辑，在这个逻辑下，一定只有一条，所以取第0个，
		user, err := user3.GetUser(i[0].Idx)
		if err != nil {
			metrics.CommonCounter.Inc("login-register", "get-userinfo-failed")
			ctx.SetError(err.Error())
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000010, Message: "获取用户信息失败"})
			return
		}
		// 生成token，无论注册还是登录均带上token返回
		token, err = tool.GenerateToken(user.Uid)
		if err != nil || token == "" {
			metrics.CommonCounter.Inc("login-register", "generate-token")
			if err != nil {
				ctx.SetError(err.Error())
			}
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000008, Message: "生成token失败"})
			return
		}
		m := user3.UserToMap(user)
		delete(m, "uid")
		result := map[string]interface{}{"token": token, "user": m}
		switch p.Method {
		case "1":
			metrics.CommonCounter.Inc("login-register", "login-1")
			// 走到这里，验证码已匹配，直接更新登录时间
			if err = user3.UpdateLoginTime(strconv.FormatInt(user.Uid, 10)); err != nil {
				metrics.CommonCounter.Inc("login-register", "update-login-time-failed-1")
				ctx.SetError(err.Error())
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000011, Message: "更新登录时间失败"}) //免密方式，更新登录时间失败
				return
			}
			metrics.CommonCounter.Inc("login-register", "login-succ-1")
			ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000002, Message: "免密方式，登录成功", Result: result})
			return
		case "2":
			metrics.CommonCounter.Inc("login-register", "login-2")
			if user.Password == "" {
				metrics.CommonCounter.Inc("login-register", "login-empty-password-2")
				// 用户通过验证码注册的，从而未设置密码(数据库中密码为空)，而登录的时候走了密码登录
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000012, Message: "您未设置密码，请使用免密登录后设置密码"})
				return
			}
			if user.Password == tool.ToMd5(p.Password) {
				// 登录成功更新登录时间
				if err = user3.UpdateLoginTime(strconv.FormatInt(user.Uid, 10)); err != nil {
					metrics.CommonCounter.Inc("login-register", "update-login-time-failed-2")
					ctx.SetError(err.Error())
					ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000013, Message: "更新登录时间失败"}) //账密方式，更新登录时间失败
					return
				}
				metrics.CommonCounter.Inc("login-register", "login-succ-2")
				ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000003, Message: "账密方式，登录成功", Result: result})
				return
			} else {
				// 用户输入的密码有误
				metrics.CommonCounter.Inc("login-register", "email-password-error")
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000014, Message: "邮箱或者密码错误"})
				return
			}
		default:
			metrics.CommonCounter.Inc("login-register", "unsupported-method")
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000007, Message: "不支持的登录/注册方式"})
			return
		}
	}
}
