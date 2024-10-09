package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
)

// ArticleDetail 文章详情页使用，通过aid获取文章信息
// 4000160
// 2000160

type articleDetailParams struct {
	Aid string `json:"aid" form:"aid" binding:"required"` // 文章id
}

func (s *Server) ArticleDetail(ctx *gin.Context) {
	params := &articleDetailParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000160, Message: "参数错误"})
		return
	}
	articleDetail, err := article.Get(params.Aid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000161, Message: "获取文章详情失败"})
		return
	}
	if err = article.IncrReadNum(params.Aid); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000162, Message: "文章阅读数+1失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000160, Message: "获取文章详情成功", Result: articleDetail})
	return
}
