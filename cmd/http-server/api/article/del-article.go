package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/article"
	"prettyy-server-online/utils/metrics"
	"strconv"
)

// DelArticle 删除文章
// 4000320
// 2000320

type delArticleParams struct {
	Aid string `form:"aid" binding:"required"` // 文章ID
	Uid int64  `form:"uid" binding:"required"` // 用户ID
}

func (s *Server) DelArticle(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("del-article", "total")
	params := &delArticleParams{}
	if err := ctx.Bind(params); err != nil {
		metrics.CommonCounter.Inc("del-article", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000320, Message: "参数错误"})
		return
	}
	ctx.SetAid(params.Aid).SetUid(strconv.FormatInt(params.Uid, 10))
	if err := article.Delete(params.Aid, params.Uid); err != nil {
		metrics.CommonCounter.Inc("del-article", "del-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000321, Message: "删除文章失败"})
		return
	}
	metrics.CommonCounter.Inc("del-article", "succ")
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000320, Message: "删除文章成功"})
	return
}
