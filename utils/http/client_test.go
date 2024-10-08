package http

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			promhttp.Handler().ServeHTTP(w, r)
			return
		}
		_, _ = w.Write([]byte("hello"))
	}))

	defer s.Close()

	client := NewClient()
	resp, _ := client.R().Get(s.URL)
	assert.Equal(t, "hello", string(resp.Body()))
	resp, _ = client.R().Get(s.URL + "/metrics")
	t.Logf("%s", string(resp.Body()))
	assert.Contains(t, string(resp.Body()), "sso_device_http_histogram_bucket")
}
