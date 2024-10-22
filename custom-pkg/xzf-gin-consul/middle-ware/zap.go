package middle_ware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
		logFields = appendLogFields("error", myCtx.GetError(), logFields)
		logFields = appendLogFields("aid", myCtx.GetAid(), logFields)
		if len(ctx.Errors) > 0 {
			logFields = append(logFields, zap.String("req_error", ctx.Errors.String()))
		}
		logger.Info("completed", logFields...)
	}
}

func appendLogFields(fieldName, fieldValue string, logFields []zap.Field) []zap.Field {
	if fieldValue == "" {
		return logFields
	}
	return append(logFields, zap.String(fieldName, fieldValue))
}

// buildLogger 构建日志，
// 开发环境同时向日志文件和屏幕均输出日志，输出所有日志级别
// 生成环境只输出到日志文件，输出Info、Warn、Error 和 DPanic、Panic、Fatal日志级别
func buildLogger() *zap.Logger {
	logMode := zapcore.InfoLevel
	if os.Getenv("idc") == "dev" {
		logMode = zapcore.DebugLevel
		core := zapcore.NewCore(getEncoder(), zapcore.NewMultiWriteSyncer(getWriteSyncer(), zapcore.AddSync(os.Stdout)), logMode)
		return zap.New(core)
	}
	core := zapcore.NewCore(getEncoder(), getWriteSyncer(), logMode)
	return zap.New(core)
}

func getWriteSyncer() zapcore.WriteSyncer {
	stSeparator := string(os.PathSeparator)
	stRootDir, _ := os.Getwd()
	stLogFilePath := stRootDir + stSeparator + "log" + stSeparator + time.Now().Format(time.DateOnly) + ".log"
	lumberjackSyncer := &lumberjack.Logger{
		Filename:   stLogFilePath,
		MaxSize:    1,     // 单位MB
		MaxBackups: 3,     // 旧文件的最大个数
		MaxAge:     30,    // 最大保存天数
		Compress:   false, // 是否压缩，disabled by default
	}
	return zapcore.AddSync(lumberjackSyncer)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		// 使用自定义的时间格式
		encoder.AppendString(t.Format(time.DateTime))
	}
	encoderConfig.CallerKey = ""     // 不记录调用者信息，删除则可自动记录行号及调用信息
	encoderConfig.StacktraceKey = "" // 不记录堆栈信息，删除则可自动记录堆栈信息
	return zapcore.NewJSONEncoder(encoderConfig)
}
