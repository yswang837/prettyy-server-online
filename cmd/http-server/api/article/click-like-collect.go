package article

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	invertedIndex "prettyy-server-online/data/inverted-index"
	"prettyy-server-online/services/article"
	invertedIndex2 "prettyy-server-online/services/inverted-index"
	"prettyy-server-online/utils/metrics"
)

// Like 点赞文章
// 2000400
// 4000400

type clickLikeCollectParams struct {
	MUid string `form:"m_uid" binding:"required"` // master uid
	SUid string `form:"s_uid" binding:"required"` // slave uid
	Aid  string `form:"aid" binding:"required"`   // 文章id
	Typ  string `form:"typ" binding:"required"`   // 4:点赞 5:收藏
}

type clickLikeCollectResp struct {
	LikeNum    int `json:"like_num"`
	CollectNum int `json:"collect_num"`
}

func (s *Server) ClickLikeCollect(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("like-collect", "total")
	var err error
	defer func() {
		if err == nil {
			metrics.CommonCounter.Inc("like-collect", "succ")
		}
	}()
	params := &clickLikeCollectParams{}
	if err = ctx.Bind(params); err != nil {
		metrics.CommonCounter.Inc("like-collect", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000400, Message: "参数错误"})
		return
	}
	ctx.SetMuid(params.MUid).SetSuid(params.SUid).SetAid(params.Aid).SetTyp(params.Typ)
	if params.Typ != invertedIndex.TypMuidLikeSuidAid && params.Typ != invertedIndex.TypMuidCollectSuidAid {
		metrics.CommonCounter.Inc("like-collect", "typ-error")
		ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000401, Message: "typ类型只能为4或5"})
		return
	}
	isLikeType := params.Typ == invertedIndex.TypMuidLikeSuidAid
	attrValue := params.MUid + "," + params.Aid
	resp := &clickLikeCollectResp{}
	if invertedIndex2.IsExist(params.Typ, attrValue, params.SUid) {
		// 说明要取消点赞或收藏，删除反向索引，并维护文章表的点赞数或者收藏数
		if err = invertedIndex2.Delete(params.Typ, attrValue, params.SUid); err != nil {
			metrics.CommonCounter.Inc("like-collect", "delete-error")
			ctx.SetError(err.Error())
			ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000402, Message: "删除反向索引失败"})
			return
		}
		if isLikeType {
			// 取消点赞类型
			likeNum, err := article.UpdateLikeNum(params.Aid, false)
			if err != nil {
				metrics.CommonCounter.Inc("like-collect", "update-like-num-error")
				ctx.SetError(err.Error())
				ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000403, Message: "更新点赞数失败"})
				return
			}
			metrics.CommonCounter.Inc("like-collect", "like-num-succ")
			resp.LikeNum = likeNum
			ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000400, Message: "取消点赞成功", Result: resp})
			return
		} else {
			// 取消收藏类型
			collectNum, err := article.UpdateCollectNum(params.Aid, false)
			if err != nil {
				metrics.CommonCounter.Inc("like-collect", "update-collect-num-error")
				ctx.SetError(err.Error())
				ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000404, Message: "更新收藏数失败"})
				return
			}
			metrics.CommonCounter.Inc("like-collect", "collect-num-succ")
			resp.CollectNum = collectNum
			ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000401, Message: "取消收藏成功", Result: resp})
			return
		}
	} else {
		// 说明要点赞或收藏，添加反向索引，并维护文章表的点赞数或者收藏数
		i := &invertedIndex.InvertedIndex{Typ: params.Typ, AttrValue: attrValue, Idx: params.SUid}
		if err = invertedIndex2.Add(i); err != nil {
			metrics.CommonCounter.Inc("like-collect", "add-error")
			ctx.SetError(err.Error())
			ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000405, Message: "添加反向索引失败"})
			return
		}
		if isLikeType {
			// 点赞类型
			likeNum, err := article.UpdateLikeNum(params.Aid, true)
			if err != nil {
				metrics.CommonCounter.Inc("like-collect", "update-like-num-error")
				ctx.SetError(err.Error())
				ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000406, Message: "更新点赞数失败"})
				return
			}
			metrics.CommonCounter.Inc("like-collect", "like-num-succ")
			resp.LikeNum = likeNum
			ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000402, Message: "点赞成功", Result: resp})
			return
		} else {
			// 收藏类型
			collectNum, err := article.UpdateCollectNum(params.Aid, true)
			if err != nil {
				metrics.CommonCounter.Inc("like-collect", "update-collect-num-error")
				ctx.SetError(err.Error())
				ctx.JSON(http.StatusBadRequest, &ginConsulRegister.Response{Code: 4000407, Message: "更新收藏数失败"})
				return
			}
			metrics.CommonCounter.Inc("like-collect", "collect-num-succ")
			resp.CollectNum = collectNum
			ctx.JSON(http.StatusOK, &ginConsulRegister.Response{Code: 2000403, Message: "收藏成功", Result: resp})
			return
		}
	}
}
