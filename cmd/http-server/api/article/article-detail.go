package article

import (
	"github.com/gin-gonic/gin"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
)

// PublishArticle 发表文章
// 4000160
// 2000160

type articleDetailParams struct {
	Aid string `json:"aid" form:"aid" binding:"required"` // 文章id
}

func (s *Server) ArticleDetail(ctx *gin.Context) {
	params := &articleDetailParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000160, Message: "参数绑定错误"})
		return
	}
	a, err := article.Get(params.Aid)
	if err != nil {
		ctx.JSON(400, ginConsulRegister.Response{Code: 4000161, Message: "获取文章详情失败"})
		return
	}
	ctx.JSON(200, ginConsulRegister.Response{Code: 2000160, Message: "获取文章详情成功", Result: a})
	return
}
