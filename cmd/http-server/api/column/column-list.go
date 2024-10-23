package column

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	invertedIndex "prettyy-server-online/data/inverted-index"
	"prettyy-server-online/services/column"
	invertedIndex2 "prettyy-server-online/services/inverted-index"
	"prettyy-server-online/utils/metrics"
	"strconv"
	"strings"
)

// ArticleList 获取专栏列表，文章详情页和专栏管理页使用
// 4000360
// 2000360

type columnListParams struct {
	Uid int64 `json:"uid" form:"uid" binding:"required"`
}

func (s *Server) ColumnList(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("column-list", "total")
	params := &columnListParams{}
	if err := ctx.Bind(params); err != nil {
		metrics.CommonCounter.Inc("column-list", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000360, Message: "参数错误"})
		return
	}
	uid := strconv.FormatInt(params.Uid, 10)
	ctx.SetUid(uid)
	cidInvertedList, err := invertedIndex2.Get(invertedIndex.TypUidCid, uid)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			metrics.CommonCounter.Inc("column-list", "inverted-column-not-found")
			ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000361, Message: "未查询到专栏数据"})
			return
		}
		metrics.CommonCounter.Inc("column-list", "get-inverted-column-err")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000362, Message: "获取反向索引失败"})
		return
	}
	result := make(map[string]string)
	for _, cidInverted := range cidInvertedList {
		col, err := column.Get(cidInverted.Idx)
		if err != nil {
			metrics.CommonCounter.Inc("column-list", "get-column-fail")
			ctx.SetError(err.Error())
			ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000363, Message: "获取专栏数据失败"})
			return
		}
		result[col.Cid] = col.Title
	}
	if len(result) == 0 {
		metrics.CommonCounter.Inc("column-list", "empty-column")
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000364, Message: "未查询到专栏数据"})
		return
	}
	metrics.CommonCounter.Inc("column-list", "succ")
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000360, Message: "获取专栏列表成功", Result: result})
	return
}
