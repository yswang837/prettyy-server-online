package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	article2 "prettyy-server-online/data/article"
	"prettyy-server-online/services/article"
	"strings"
)

// ArticleList 获取文章列表
// 4000180
// 2000180

type articleListParams struct {
	Page     int `json:"page" form:"page" binding:"required"` // 第几页
	PageSize int `json:"page_size" form:"page_size"`          // 每页多少条
}

func (s *Server) ArticleList(ctx *gin.Context) {
	params := &articleListParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000180, Message: "参数错误"})
		return
	}
	a, err := article.GetArticleList(params.Page, params.PageSize)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000181, Message: "没有更多数据", Result: []article2.Article{}})
			return
		}
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000182, Message: "获取文章列表失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000180, Message: "获取文章列表成功", Result: a})
	return
}
