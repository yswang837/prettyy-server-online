package middle_ware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func NewZapLogger() func(ctx *gin.Context) {
	logger, _ := zap.NewProduction()
	return func(ctx *gin.Context) {
		start := time.Now()
		// 处理请求
		ctx.Next()
		endTime := time.Since(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		url := ctx.Request.RequestURI
		status := ctx.Writer.Status()
		clientUserAgent := ctx.Request.UserAgent()
		clientProtocol := ctx.Request.Proto
		logFields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("url", url),
			zap.String("client_ip", clientIP),
			zap.String("client_user_agent", clientUserAgent),
			zap.String("client_protocol", clientProtocol),
			zap.String("exec_time", strconv.Itoa(int(endTime.Milliseconds()))+"ms"),
		}
		if len(ctx.Errors) > 0 {
			logFields = append(logFields, zap.String("error", ctx.Errors.String()))
			logger.Error("error occurred during request", logFields...)
		} else {
			logger.Info("request completed", logFields...)
		}
	}
}
