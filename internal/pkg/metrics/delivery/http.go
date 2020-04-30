package delivery

import (
	"failless/internal/pkg/metrics/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var promHandler = promhttp.Handler()

func MetricsHandler(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	promHandler.ServeHTTP(w, r)
}

//NewPrometheusService create a new prometheus service
func NewPrometheusService() (*usecase.MetricUseCase, error) {
	cli := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "pushgateway",
		Name:      "cmd_duration_seconds",
		Help:      "Client application execution in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"name"})
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &usecase.MetricUseCase{
		PHistogram:           cli,
		HttpRequestHistogram: http,
	}
	err := prometheus.Register(s.PHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	err = prometheus.Register(s.HttpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}
