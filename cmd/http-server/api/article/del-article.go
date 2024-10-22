package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
)

// DelArticle 删除文章
// 4000320
// 2000320

type delArticleParams struct {
	Aid string `json:"aid" form:"aid" binding:"required"` // 文章ID
	Uid int64  `json:"uid" form:"uid" binding:"required"` // 用户ID
}

func (s *Server) DelArticle(ctx *ginConsulRegister.Context) {
	params := &delArticleParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000320, Message: "参数错误"})
		return
	}
	if err := article.Delete(params.Aid, params.Uid); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000321, Message: "删除文章失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000320, Message: "删除文章成功"})
	return
}
