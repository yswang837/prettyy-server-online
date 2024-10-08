package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"prettyy-server-online/utils/tool"
)

type counter struct {
	promCounterVec *prometheus.CounterVec
}

type histogram struct {
	promHistogramVec *prometheus.HistogramVec
}

type gauge struct {
	proGaugeVec *prometheus.GaugeVec
}

func (c *counter) Values(v ...string) prometheus.Counter {
	v = append([]string{tool.ServerName}, v...)
	return c.promCounterVec.WithLabelValues(v...)
}

func newCounter(name, help string, tags []string) *counter {
	tags = append([]string{"service_name"}, tags...)
	vec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: os.Getenv("NAMESPACE"),
		Subsystem: os.Getenv("SUBSYSTEM"),
		Name:      tool.ServerName + "_" + name,
		Help:      help,
	}, tags)
	c := &counter{vec}
	prometheus.MustRegister(vec)
	return c
}

func newHistogram(name, help string, tags []string, buckets []float64) *histogram {
	tags = append([]string{"service_name"}, tags...)
	vec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      tool.ServerName + "_" + name,
			Namespace: os.Getenv("NAMESPACE"),
			Help:      help,
			Buckets:   buckets,
		}, tags)
	h := &histogram{vec}
	prometheus.MustRegister(vec)
	return h
}

func (h *histogram) Values(v ...string) prometheus.Observer {
	v = append([]string{tool.ServerName}, v...)
	return h.promHistogramVec.WithLabelValues(v...)
}

func newGauge(name string, help string, tags []string) *gauge {
	tags = append([]string{"service_name"}, tags...)
	vec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: os.Getenv("NAMESPACE"),
			Subsystem: os.Getenv("SUBSYSTEM"),
			Name:      tool.ServerName + "_" + name,
			Help:      help,
		}, tags)
	prometheus.MustRegister(vec)
	g := &gauge{vec}
	return g
}

func (g *gauge) Values(v ...string) prometheus.Gauge {
	v = append([]string{tool.ServerName}, v...)
	return g.proGaugeVec.WithLabelValues(v...)
}
