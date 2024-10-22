package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	article2 "prettyy-server-online/data/article"
	"prettyy-server-online/services/article"
	"strings"
)

// ArticleList 获取文章列表，如果uid为小于1w，则获取所有文章供首页使用，否则获取该用户的所有文章，支持分页
// 4000180
// 2000180

type articleListParams struct {
	Uid        int64  `json:"uid" form:"uid"`                      // uid 如果合法，则返回的是当前用户的文章列表，否则返回的是首页文章的列表
	Page       int    `json:"page" form:"page" binding:"required"` // 第几页
	PageSize   int    `json:"page_size" form:"page_size"`          // 每页多少条
	Visibility string `json:"visibility" form:"visibility"`        // 文章的可见性, 1-全部可见 2-VIP可见 3-粉丝可见 4-仅我可见
	Typ        string `json:"typ" form:"typ"`                      // 文章类型，1-原创 2-转载 3-翻译
}

func (s *Server) ArticleList(ctx *ginConsulRegister.Context) {
	params := &articleListParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000180, Message: "参数错误"})
		return
	}
	a, count, err := article.GetArticleList(params.Uid, params.Page, params.PageSize, params.Visibility, params.Typ)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000181, Message: "没有更多数据", Result: []article2.Article{}})
			return
		}
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000182, Message: "获取文章列表失败"})
		return
	}
	result := map[string]interface{}{"article_list": a, "count": count}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000180, Message: "获取文章列表成功", Result: result})
	return
}
