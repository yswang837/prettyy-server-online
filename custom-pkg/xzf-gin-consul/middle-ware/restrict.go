package middle_ware

import "github.com/gin-gonic/gin"

// Restrict 限制访问，目前是针对邮件发送验证码来做，1分钟之内发送1次，1小时之内发送5次，一天之内发送20次
func Restrict() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
