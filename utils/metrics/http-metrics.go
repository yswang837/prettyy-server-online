package metrics

var HttpHistogramMetrics = newHttpHistogramMetrics()

// HttpHistogramTags ...
type HttpHistogramTags struct {
	Domain string // Http host
	Method string // 操作类型
	Status string // 操作类型
	URL    string // Http实例名称
}

type httpHistogramMetrics struct {
	histogram *histogram
}

func (acc *httpHistogramMetrics) Observe(tags *HttpHistogramTags, f float64) {
	acc.histogram.Values(tags.Domain, tags.Method, tags.Status, tags.URL).Observe(f)
}

func newHttpHistogramMetrics() *httpHistogramMetrics {
	acc := &httpHistogramMetrics{
		histogram: newHistogram("http_histogram", "Http query time histogram", []string{"domain", "method", "status", "url"}, []float64{2, 5, 10, 25, 50, 100}),
	}
	return acc
}
