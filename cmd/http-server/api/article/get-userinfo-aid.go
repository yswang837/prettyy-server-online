package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
	user3 "prettyy-server-online/services/user"
	"strconv"
)

// GetUserInfoByAid 文章详情页使用，通过aid获取用户信息
// 2000340
// 4000340

type getUserInfoByAidParams struct {
	Aid string `json:"aid" form:"aid" binding:"required"` // 文章id
}

func (s *Server) GetUserInfoByAid(ctx *ginConsulRegister.Context) {
	params := &getUserInfoByAidParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000340, Message: "参数错误"})
		return
	}
	a, err := article.Get(params.Aid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000341, Message: "获取文章信息失败"})
		return
	}
	u, err := user3.GetUser(strconv.FormatInt(a.Uid, 10))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000342, Message: "获取用户信息失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000340, Message: "获取用户信息成功", Result: u})
	return
}
