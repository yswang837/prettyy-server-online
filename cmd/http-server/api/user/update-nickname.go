package user

import (
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
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
	p := &updateNickNameParams{}
	if err := ctx.Bind(p); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000200, Message: "参数错误"})
		return
	}
	u, err := user.GetUser(p.Uid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000201, Message: "获取用户信息失败"})
		return
	}
	if p.NickName == u.NickName {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000202, Message: "昵称未改变"})
		return
	}
	if err = user.UpdateNickName(p.Uid, p.NickName); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000203, Message: "更新昵称失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000200, Message: "更新昵称成功"})
	return
}
