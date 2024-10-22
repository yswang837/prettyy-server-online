package middle_ware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"strconv"
	"time"
)

func NewZapLogger() func(ctx *gin.Context) {
	logger := buildLogger()
	return func(ctx *gin.Context) {
		start := time.Now()
		myCtx := ginConsulRegister.NewContext(ctx)
		// 处理请求
		ctx.Next()
		endTime := time.Since(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		url := ctx.Request.RequestURI
		status := ctx.Writer.Status()
		clientUserAgent := ctx.Request.UserAgent()
		clientProtocol := ctx.Request.Proto
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		logFields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("url", url),
			zap.String("client_ip", clientIP),
			zap.String("client_user_agent", clientUserAgent),
			zap.String("client_protocol", clientProtocol),
			zap.String("exec_time", strconv.Itoa(int(endTime.Milliseconds()))+"ms"),
			zap.String("hostname", hostName),
			zap.String("caller", myCtx.GetCaller()),
		}
		if len(ctx.Errors) > 0 {
			logFields = append(logFields, zap.String("req_error", ctx.Errors.String()))
			logger.Error("request error", logFields...)
		} else {
			myCtxErr := myCtx.GetError()
			if myCtxErr != "" {
				logFields = append(logFields, zap.String("sys_error", myCtxErr))
				logger.Error("system error", logFields...)
			} else {
				logger.Info("request completed", logFields...)
			}
		}
	}
}

func buildLogger() *zap.Logger {
	// 获取生产环境的配置
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		// 使用自定义的时间格式
		encoder.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	config.EncoderConfig.CallerKey = ""     // 不记录调用者信息，删除则可自动记录行号及调用信息
	config.EncoderConfig.StacktraceKey = "" // 不记录堆栈信息，删除则可自动记录堆栈信息
	// 创建 logger
	log, _ := config.Build()
	return log
}
