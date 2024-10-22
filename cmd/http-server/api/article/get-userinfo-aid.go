package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
	user3 "prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
	"strconv"
)

// GetUserInfoByAid 文章详情页使用，通过aid获取用户信息
// 2000340
// 4000340

type getUserInfoByAidParams struct {
	Aid string `form:"aid" binding:"required"` // 文章id
}

func (s *Server) GetUserInfoByAid(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("get-userinfo-aid", "total")
	params := &getUserInfoByAidParams{}
	if err := ctx.Bind(params); err != nil {
		metrics.CommonCounter.Inc("get-userinfo-aid", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000340, Message: "参数错误"})
		return
	}
	ctx.SetAid(params.Aid)
	a, err := article.Get(params.Aid)
	if err != nil {
		metrics.CommonCounter.Inc("get-userinfo-aid", "get-articleinfo-fail")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000341, Message: "获取文章信息失败"})
		return
	}
	u, err := user3.GetUser(strconv.FormatInt(a.Uid, 10))
	if err != nil {
		metrics.CommonCounter.Inc("get-userinfo-aid", "get-userinfo-fail")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000342, Message: "获取用户信息失败"})
		return
	}
	metrics.CommonCounter.Inc("get-userinfo-aid", "succ")
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000340, Message: "获取用户信息成功", Result: u})
	return
}
