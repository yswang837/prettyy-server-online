package http

import (
	"github.com/go-resty/resty/v2"
	"prettyy-server-online/utils/metrics"
	"strconv"
	"time"
)

func NewClient(handler ...Handler) *resty.Client {
	client := resty.New()
	client.SetTimeout(time.Second)
	client.OnAfterResponse(MetricsMiddleware(handler...)) //注册插件
	return client
}

// MetricsMiddleware to measure time of http request
func MetricsMiddleware(handler ...Handler) resty.ResponseMiddleware {
	return func(_ *resty.Client, response *resty.Response) error {
		for _, h := range handler {
			_ = h(response.Request, response)
		}
		duration := response.Time().Milliseconds() // 不包含执行Middleware的时间
		req := response.Request.RawRequest
		metrics.HttpHistogramMetrics.Observe(&metrics.HttpHistogramTags{
			Domain: req.Host,
			Method: req.Method,
			Status: strconv.Itoa(response.StatusCode()),
			URL:    req.URL.Path, // 注意： 某些restful的请求不适合测量
		}, float64(duration))
		return nil
	}
}

type Handler func(req *resty.Request, resp *resty.Response) error
