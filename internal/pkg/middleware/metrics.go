package middleware

import (
	"failless/internal/pkg/metrics"
	"failless/internal/pkg/metrics/delivery"
	"failless/internal/pkg/settings"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var mUseCase metrics.UseCase
var doOnce sync.Once

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func Metrics(next settings.HandlerFunc) settings.HandlerFunc {
	doOnce.Do(func() {
		var err error = nil
		mUseCase, err = delivery.NewPrometheusService()
		if err != nil {
			log.Fatal("METRICS MIDDLEWARE ERROR: ", err.Error())
		}
	})
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		appMetric := metrics.NewHTTP(r.URL.Path, r.Method)
		appMetric.Started()
		sw := statusWriter{ResponseWriter: w}
		next(&sw, r, ps)
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(sw.status)
		mUseCase.SendMetric(appMetric)
	}
}
