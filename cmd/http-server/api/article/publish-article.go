package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	xzfSnowflake "prettyy-server-online/custom-pkg/xzf-snowflake"
	article2 "prettyy-server-online/data/article"
	invertedIndex "prettyy-server-online/data/inverted-index"
	"prettyy-server-online/services/article"
	"prettyy-server-online/services/column"
	invertedIndex2 "prettyy-server-online/services/inverted-index"
	"prettyy-server-online/utils/tool"
	"strconv"
	"strings"
)

// PublishArticle 发表文章
// 4000120
// 2000120

const (
	indexAidTyp   = "2" // uid->aid
	indexCidTyp   = "3"
	articlePrefix = "AA"
	columnPrefix  = "AB"
)

type articleParams struct {
	Title      string `json:"title" form:"title" binding:"required"`         // 文章标题
	Content    string `json:"content" form:"content" binding:"required"`     // 文章内容
	CoverImg   string `json:"cover_img" form:"cover_img" binding:"required"` // 文章封面url
	Summary    string `json:"summary" form:"summary" binding:"required"`     // 文章摘要
	Visibility string `json:"visibility" form:"visibility"`                  // 文章的可见性，默认全部可见 "1"-全部可见 "2"-VIP可见 "3"-粉丝可见 "4"-仅我可见
	Tags       string `json:"tags" form:"tags" binding:"required"`           // 文章标签，以英文逗号分隔，最多10个标签，由用户发文的时候打标签
	Column     string `json:"column" form:"column"`                          // 文章所属专栏，默认为空，非空则以英文逗号分隔，最多4个专栏，由用户发文的时候指定，格式：cid1,title1,"",title2,cid3,title3
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
		Aid:        xzfSnowflake.GenID(articlePrefix),
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
	// 先查反向索引，查到了就追加aid，查不到就添加aid
	// 向反向索引中添加uid->aid1,aid2,...的映射，以便在开启分表后，内容管理页面查询当前用户的文章
	uid := strconv.FormatInt(params.Uid, 10)
	idba, _ := invertedIndex2.Get(indexAidTyp, uid)
	if idba == nil {
		i := &invertedIndex.InvertedIndex{Typ: indexAidTyp, AttrValue: uid, Index: a.Aid}
		if err := invertedIndex2.Add(i); err != nil {
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000122, Message: "添加文章的反向索引失败"})
			return
		}
	} else {
		// 在反向索引的index字段中，设置的最大长度是varchar(8192),一个用户最多430篇文章
		idba.Index = idba.Index + "," + a.Aid
		if err := invertedIndex2.UpdateAid(indexAidTyp, uid, idba.Index); err != nil {
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000123, Message: "更新文章的反向索引失败"})
			return
		}
	}
	// params.Column格式：cid1,title1,"",title2,cid3,title3
	// 从params.Column中获取，如果cid存在，则跳过该cid和title，如果不存在则将其写入到needInsertToColumn
	// needInsertToColumn 生成新的cid:title，用于写入专栏表和更新反向索引表
	needInsertToColumn := make(map[string]string)
	if params.Column != "" {
		columns := strings.Split(params.Column, ",")
		length := len(columns)
		if length == 0 || length >= 9 || length%2 != 0 {
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000124, Message: "专栏数量不合法"})
			return
		}
		for i := 0; i < length; i = i + 2 {
			if len(columns[i]) == 19 && strings.HasPrefix(columns[i], columnPrefix) {
				continue
			}
			needInsertToColumn[xzfSnowflake.GenID(columnPrefix)] = columns[i+1]
		}
	}
	if len(needInsertToColumn) != 0 {
		// 维护专栏表
		if err := column.Add(needInsertToColumn, params.Uid); err != nil {
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000125, Message: "添加专栏失败"})
		}
		// 维护反正索引
		idbc, _ := invertedIndex2.Get(indexCidTyp, uid)
		var cidSlice []string
		for cid := range needInsertToColumn {
			cidSlice = append(cidSlice, cid)
		}
		cidStr := strings.Join(cidSlice, ",")
		if idbc == nil {
			i := &invertedIndex.InvertedIndex{Typ: indexCidTyp, AttrValue: uid, Index: cidStr}
			if err := invertedIndex2.Add(i); err != nil {
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000126, Message: "添加专栏的反向索引失败"})
				return
			}
		} else {
			// 在反向索引的index字段中，设置的最大长度是varchar(8192),一个用户最多430个专栏
			idbc.Index = idbc.Index + "," + cidStr
			if err := invertedIndex2.UpdateCid(indexCidTyp, uid, idbc.Index); err != nil {
				ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000127, Message: "更新专栏的反向索引失败"})
				return
			}
		}
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000120, Message: "添加文章成功"})
	return
}
