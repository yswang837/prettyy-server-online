package middle_ware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	ginConsulRegister "prettyy-server-online/custom-pkg/xzf-gin-consul/register"
	"strconv"
	"time"
)

func NewMetrics(prefix string) func(ctx *gin.Context) {
	tags := []string{"service_name", "domain", "idc", "caller", "url", "http_status", "message"}
	buckets := []float64{2, 5, 10, 25, 50, 100}
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    prefix + "http_request_histogram",
			Buckets: buckets,
		}, tags)
	prometheus.MustRegister(vec)

	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			domain := ctx.Request.Host
			url := ctx.Request.URL.Path
			status := ctx.Writer.Status()
			if status == 404 || status == 403 {
				domain = "unknown"
				url = "unknown"
			}
			myCtx := ginConsulRegister.NewContext(ctx)
			message := myCtx.GetMessage()
			caller := myCtx.GetCaller()

			vec.WithLabelValues(os.Getenv("SERVICE_NAME"), domain, os.Getenv("IDC"), caller, url, strconv.Itoa(status), message).Observe(float64(time.Since(start) / time.Millisecond))
		}()
		ctx.Next()
	}
}

type MetricService struct {
	uri string
}

func NewMetricService() *MetricService {
	return &MetricService{uri: "/metrics"}
}

func (m *MetricService) Init() error {
	return nil
}

func (m *MetricService) SetRoute(r *gin.Engine) {
	r.GET(m.uri, gin.WrapH(promhttp.Handler()))
}
