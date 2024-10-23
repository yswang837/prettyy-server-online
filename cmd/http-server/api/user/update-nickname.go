package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/metrics"
)

// updateNickNameParams 面向接口
type updateNickNameParams struct {
	Uid      string `json:"uid" form:"uid" binding:"required"`
	NickName string `json:"nick_name" form:"nick_name" binding:"required"`
}

// UpdateNickName 更新用户名，默认用户名为邮箱前缀
// 4000200
// 2000200
func (s *Server) UpdateNickName(ctx *ginConsulRegister.Context) {
	metrics.CommonCounter.Inc("update-nickname", "total")
	p := &updateNickNameParams{}
	if err := ctx.Bind(p); err != nil {
		metrics.CommonCounter.Inc("update-nickname", "params-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000200, Message: "参数错误"})
		return
	}
	ctx.SetUid(p.Uid).SetNickname(p.NickName)
	u, err := user.GetUser(p.Uid)
	if err != nil {
		metrics.CommonCounter.Inc("update-nickname", "get-user-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000201, Message: "获取用户信息失败"})
		return
	}
	if p.NickName == u.NickName {
		metrics.CommonCounter.Inc("update-nickname", "same-nickname")
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000202, Message: "昵称未改变"})
		return
	}
	if err = user.UpdateNickName(p.Uid, p.NickName); err != nil {
		metrics.CommonCounter.Inc("update-nickname", "update-nickname-error")
		ctx.SetError(err.Error())
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000203, Message: "更新昵称失败"})
		return
	}
	metrics.CommonCounter.Inc("update-nickname", "succ")
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000200, Message: "更新昵称成功"})
	return
}
