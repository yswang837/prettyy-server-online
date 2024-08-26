package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	article2 "prettyy-server-online/data/article"
	"prettyy-server-online/services/article"
	"prettyy-server-online/utils/tool"
)

// PublishArticle 发表文章
// 4000120
// 2000120

type articleParams struct {
	Title    string `json:"title" form:"title" binding:"required"`         // 文章标题
	Content  string `json:"content" form:"content" binding:"required"`     // 文章内容
	Uid      int64  `json:"uid" form:"uid" binding:"required"`             // 用户id
	CoverImg string `json:"cover_img" form:"cover_img" binding:"required"` // 文章封面url
	Summary  string `json:"summary" form:"summary" binding:"required"`     // 文章摘要

}

func (s *Server) PublishArticle(ctx *gin.Context) {
	params := &articleParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000120, Message: "参数错误"})
		return
	}
	a := &article2.Article{
		Title:    params.Title,
		Content:  tool.Base64Encode(params.Content),
		CoverImg: params.CoverImg,
		Summary:  params.Summary,
		Uid:      params.Uid,
	}
	if err := article.Add(a); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000121, Message: "添加文章失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000120, Message: "添加文章成功"})
	return
}
