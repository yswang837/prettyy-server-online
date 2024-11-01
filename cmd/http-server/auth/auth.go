package auth

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gin-gonic/gin"
	"net/url"
	"prettyy-server-online/cmd/http-server/conf"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/utils/tool"
	"sort"
	"strings"
)

var (
	authConfig = map[string]struct {
		allowCaller *hashset.Set
	}{
		conf.URLIdentifyCodeByEmail: {allowCaller: hashset.New("test", "web")},
		conf.URLIdentifyCode:        {allowCaller: hashset.New("test", "web")},
		conf.URLRegisterLogin:       {allowCaller: hashset.New("test", "web")},
		conf.URLCheckPassword:       {allowCaller: hashset.New("test", "web")},
		conf.URLLoginOut:            {allowCaller: hashset.New("test", "web")},
		conf.URLUpdateNickName:      {allowCaller: hashset.New("test", "web")},
		conf.URLUpdateGender:        {allowCaller: hashset.New("test", "web")},
		conf.URLUpdateSummary:       {allowCaller: hashset.New("test", "web")},
		conf.URLUpdateProvinceCity:  {allowCaller: hashset.New("test", "web")},
		conf.URLUpdateBirthday:      {allowCaller: hashset.New("test", "web")},
		conf.URLUpdatePassword:      {allowCaller: hashset.New("test", "web")},
		conf.URLPublishArticle:      {allowCaller: hashset.New("test", "web")},
		conf.URLGetArticleList:      {allowCaller: hashset.New("test", "web")},
		conf.URLGetArticleDetail:    {allowCaller: hashset.New("test", "web")},
		conf.URLDelArticle:          {allowCaller: hashset.New("test", "web")},
		conf.URLGetUserInfoByAid:    {allowCaller: hashset.New("test", "web")},
		conf.URLGetColumnList:       {allowCaller: hashset.New("test", "web")},
		conf.URLFileUpload:          {allowCaller: hashset.New("test", "web")},
		conf.URLExtractSummary:      {allowCaller: hashset.New("test", "web")},
		conf.URLLikeCollectArticle:  {allowCaller: hashset.New("test", "web")},
		conf.URLChat:                {allowCaller: hashset.New("test", "web")},
	}
	callerPin = map[string]string{
		"test": "098f6bcd4621d373cade4e832627b4f6", // 测试的调用方和pin，
		"web":  "4d3abacc23255a0767c718afda30cc8c", // web浏览器的调用方和pin，
	}
)

func Auth(ginCtx *gin.Context) {
	ctx := ginConsulRegister.NewContext(ginCtx)
	// 开发业务代码阶段，可以临时注释一下
	if !auth(ctx) {
		ctx.Abort()
		return
	}
	ctx.Next()
}

func auth(ctx *ginConsulRegister.Context) bool {
	// 检查url是否允许
	urlStr := ctx.Request.URL.Path
	config, ok := authConfig[urlStr]
	if !ok {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400001, Message: fmt.Sprintf("url %s is not allow", urlStr)})
		return false
	}
	// 检查caller是否允许
	caller := getParam(ctx, "caller")
	if caller == "" {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400002, Message: "caller is empty"})
		return false
	}
	ctx.SetCaller(caller)
	if config.allowCaller != nil && !config.allowCaller.Contains(caller) {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400003, Message: fmt.Sprintf("caller: %s is forbidden for %s", caller, urlStr)})
		return false
	}
	// 检查签名是否正确
	sign := getParam(ctx, "sign")
	if sign == "" {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400004, Message: "sign is empty"})
		return false
	}
	signStr := getParam(ctx, "")
	expectedSign := tool.ToMd5(callerPin[caller] + signStr)
	if expectedSign != sign {
		ctx.JSON(200, &ginConsulRegister.Response{Code: 400005, Message: fmt.Sprintf("signature error. expected %s, actuall %s", expectedSign, sign)})
		return false
	}
	return true
}

func getParam(ctx *ginConsulRegister.Context, key string) string {
	// key如果不为空，则获取单个属性值，
	// key如果为空则获取所有属性值，并且是排好序的值的字符串拼接，用于验签
	switch ctx.Request.Method {
	case "GET":
		if key != "" {
			value, _ := ctx.GetQuery(key)
			return value
		} else {
			return sortParams(ctx.Request.URL.Query())
		}
	case "POST":
		if key != "" {
			return ctx.PostForm(key)
		} else {
			// 获取所有的参数
			return sortParams(ctx.Request.PostForm)
		}
	}
	return ""
}

func sortParams(m url.Values) string {
	if len(m) == 0 {
		return ""
	}
	var params []string
	for k, v := range m {
		// 上传文件接口，file字段本身就获取不到，这里也写明直接排除
		// 发布文章接口，content内容多且杂，也排除
		if k == "sign" || k == "file" || k == "content" {
			continue
		}
		if len(v) != 0 {
			params = append(params, fmt.Sprintf("%s=%s", k, v[0]))
		}
	}
	sort.Strings(params)
	return strings.Join(params, "")
}
