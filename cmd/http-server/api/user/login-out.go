package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/tool"
)

// LoginOut 用户点击退出登录，将token加入到token的黑名单cache中，该缓存的过期时间为默认的1小时
// 4000060
// 2000060
func (s *Server) LoginOut(ctx *gin.Context) {
	token := tool.GetToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000060, Message: "token为空"})
		return
	}
	if err := user.SetExByToken(token); err != nil {
		ctx.JSON(http.StatusBadRequest, ginConsulRegister.Response{Code: 4000061, Message: "token作废失败"})
		return
	}
	ctx.JSON(http.StatusOK, ginConsulRegister.Response{Code: 2000060, Message: "token作废成功"})
	return
}
