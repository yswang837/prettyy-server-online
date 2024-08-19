package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"prettyy-server-online/cmd/http-server/conf"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
)

var (
	authConfig = map[string]struct {
		signAttr    []string
		allowCaller *hashset.Set
	}{
		conf.URLIdentifyCodeByEmail: {signAttr: []string{"email"}, allowCaller: hashset.New("test")},
		conf.URLIdentifyCode:        {signAttr: []string{""}, allowCaller: hashset.New("test")},
		conf.URLRegisterLogin:       {signAttr: []string{"email", "password", "identity_code"}, allowCaller: hashset.New("test")},
		conf.URLCheckPassword:       {signAttr: []string{"email"}, allowCaller: hashset.New("test")},
		conf.URLLoginOut:            {signAttr: []string{""}, allowCaller: hashset.New("test")},
		conf.URLUpdateNickName:      {signAttr: []string{"email", "nick_name"}, allowCaller: hashset.New("test")},
		conf.URLUpdateGender:        {signAttr: []string{"email", "gender"}, allowCaller: hashset.New("test")},
		conf.URLUpdateSummary:       {signAttr: []string{"email", "summary"}, allowCaller: hashset.New("test")},
		conf.URLUpdateProvinceCity:  {signAttr: []string{"email", "province", "city"}, allowCaller: hashset.New("test")},
		conf.URLUpdateBirthday:      {signAttr: []string{"email"}, allowCaller: hashset.New("test")},
		conf.URLPublishArticle:      {signAttr: []string{"title", "summary", "uid"}, allowCaller: hashset.New("test")},
		conf.URLGetArticleList:      {signAttr: []string{"page", "page_size"}, allowCaller: hashset.New("test")},
		conf.URLGetArticleDetail:    {signAttr: []string{"page", "page_size"}, allowCaller: hashset.New("test")},
		conf.URLFileUpload:          {signAttr: []string{""}, allowCaller: hashset.New("test")},
	}
	callerPin = map[string]string{
		"test": "098f6bcd4621d373cade4e832627b4f6", // 测试的调用方和pin，
	}
)

func Auth(ginCtx *gin.Context) {
	ctx := ginConsulRegister.NewContext(ginCtx)
	// 开发业务代码阶段，可以临时注释一下
	//if !auth(ctx) {
	//	ctx.Abort()
	//	return
	//}
	ctx.Next()
}

func auth(ctx *ginConsulRegister.Context) bool {
	// 检查url是否允许
	url := ctx.Request.URL.Path
	config, ok := authConfig[url]
	if !ok {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400001, Message: fmt.Sprintf("url %s is not allow", url)})
		return false
	}
	var getParam = func(key string) string {
		value, ok := ctx.GetQuery(key)
		if !ok {
			value = ctx.PostForm(key)
		}
		return value
	}
	// 检查caller是否允许
	caller := getParam("caller")
	if caller == "" {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400002, Message: "caller is empty"})
		return false
	}
	ctx.SetCaller(caller)
	if config.allowCaller != nil && !config.allowCaller.Contains(caller) {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400003, Message: fmt.Sprintf("caller: %s is forbidden for %s", caller, url)})
		return false
	}

	// 检查签名是否正确
	var signSlice []string
	for _, attr := range config.signAttr {
		signSlice = append(signSlice, getParam(attr))
	}
	signStr := ""
	for _, value := range signSlice {
		signStr += value
	}
	sign := getParam("sign")
	b := md5.Sum([]byte(callerPin[caller] + signStr))
	expectedSign := hex.EncodeToString(b[:])
	if expectedSign != sign {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400004, Message: fmt.Sprintf("signature error. expected %s, actuall %s", expectedSign, sign)})
		return false
	}
	return true
}
