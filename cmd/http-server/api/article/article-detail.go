package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
	"prettyy-server-online/utils/metrics"
)

// ArticleDetail 文章详情页使用，通过aid获取文章信息
// 4000160
// 2000160

type articleDetailParams struct {
	Aid string `form:"aid" binding:"required"` // 文章id
}

func (s *Server) ArticleDetail(ctx *ginConsulRegister.Context) {
	// todo 不知道为什么前端多次刷新网页，这个接口就卡住了，通过apipost多次请求一样的，基本上是后端的问题，请求是进来了的，但是掉redis和mysql就一直卡住
	// 可能是 hgetall 查文章内容卡住了 redis，先注释试试
	metrics.CommonCounter.Inc("article-detail", "total")
	params := &articleDetailParams{}
	if err := ctx.Bind(params); err != nil {
		metrics.CommonCounter.Inc("article-detail", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000160, Message: "参数错误"})
		return
	}
	ctx.SetAid(params.Aid)
	if err := article.IncrReadNum(params.Aid); err != nil {
		metrics.CommonCounter.Inc("article-detail", "incr-read-num-err")
		ctx.SetError(err.Error())
		// 这是无关紧要的失败，感觉可以注释掉
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000161, Message: "文章阅读数+1失败"})
		return
	}
	articleDetail, err := article.Get(params.Aid)
	if err != nil {
		metrics.CommonCounter.Inc("article-detail", "get-article-detail")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000162, Message: "获取文章详情失败"})
		return
	}
	metrics.CommonCounter.Inc("article-detail", "succ")
	ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000160, Message: "获取文章详情成功", Result: articleDetail})
	return
}
