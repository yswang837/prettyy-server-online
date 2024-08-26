package middle_ware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"prettyy-server-online/services/user"
	"prettyy-server-online/utils/tool"
)

// JwtAuth 中间件，检查token
// 2000080
// 4000080
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := tool.GetToken(c) // 从 header 或者 cookie 中获取 token
		if token == "" {
			c.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000080, Message: "未登录或非法访问"})
			c.Abort()
			return
		}
		// 判断是否已经加入到token的黑名单中，如果是，帐户异地登陆或令牌失效
		if user.IsExistToken(token) {
			c.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000081, Message: "帐户异地登陆或令牌失效"})
			c.Abort()
			return
		}
		// 解析 token
		j, err := tool.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, ginConsulRegister.Response{Code: 4000082, Message: "token未授权"})
			c.Abort()
			return
		}
		c.Set("Authorization", j)
		c.Next()
	}
}
