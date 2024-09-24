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
	Title      string `json:"title" form:"title" binding:"required"`         // 文章标题
	Content    string `json:"content" form:"content" binding:"required"`     // 文章内容
	CoverImg   string `json:"cover_img" form:"cover_img" binding:"required"` // 文章封面url
	Summary    string `json:"summary" form:"summary" binding:"required"`     // 文章摘要
	Visibility string `json:"visibility" form:"visibility"`                  // 文章的可见性，默认全部可见 "1"-全部可见 "2"-VIP可见 "3"-粉丝可见 "4"-仅我可见
	Tags       string `json:"tags" form:"tags" binding:"required"`           // 文章标签，以英文逗号分隔，最多10个标签，由用户发文的时候打标签
	Typ        string `json:"typ" form:"typ"`                                // 文章类型，默认原创，"1"-原创 "2"-转载 "3"-翻译 "4"-其他
	Uid        int64  `json:"uid" form:"uid" binding:"required"`             // 用户id
}

func (s *Server) PublishArticle(ctx *gin.Context) {
	params := &articleParams{}
	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000120, Message: "参数错误"})
		return
	}
	a := &article2.Article{
		Title:      params.Title,
		Content:    tool.Base64Encode(params.Content),
		CoverImg:   params.CoverImg,
		Summary:    params.Summary,
		Visibility: params.Visibility,
		Tags:       params.Tags,
		Typ:        params.Typ,
		Uid:        params.Uid,
	}
	if err := article.Add(a); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000121, Message: "添加文章失败"})
		return
	}
	// todo 向反向索引中添加uid->aid1,aid2,...的映射，以便在开启分表后，内容管理页面查询当前用户的文章
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000120, Message: "添加文章成功"})
	return
}
