package middleware

import (
	"net/http"
	"strconv"
	"time"
	
	"github.com/HubertLipinski/go-rest-graphql-grpc/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore prometheus data endpoint
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		rec := &statusRecorder{ResponseWriter: w, status: 200}

		start := time.Now()
		next.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		status := strconv.Itoa(rec.status)

		metrics.HttpRequestsTotal.WithLabelValues(path, method, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(path, method, status).Observe(duration)
	})
}
