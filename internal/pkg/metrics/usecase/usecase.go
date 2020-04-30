package usecase

import (
	"failless/internal/pkg/metrics"
	"failless/internal/pkg/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// Service implements UseCase interface
type MetricUseCase struct {
	PHistogram           *prometheus.HistogramVec
	HttpRequestHistogram *prometheus.HistogramVec
}

// SaveNewClient send metrics to server
func (s *MetricUseCase) SaveNewClient(c *metrics.Client) error {
	gatewayURL := settings.SecureSettings.MetricsHost
	s.PHistogram.WithLabelValues(c.Name).Observe(c.Duration)
	return push.New(gatewayURL, "cmd_job").Collector(s.PHistogram).Push()
}

// SendMetric send metrics to server
func (s *MetricUseCase) SendMetric(h *metrics.HTTP) {
	s.HttpRequestHistogram.WithLabelValues(h.Handler, h.Method, h.StatusCode).Observe(h.Duration)
}
