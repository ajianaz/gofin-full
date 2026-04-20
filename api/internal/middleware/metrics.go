package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gofin_http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gofin_http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"method", "path"})

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gofin_http_requests_in_flight",
		Help: "Number of HTTP requests currently being processed",
	})

	httpResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "gofin_http_response_size_bytes",
		Help:    "HTTP response size in bytes",
		Buckets: prometheus.ExponentialBuckets(100, 10, 7),
	}, []string{"method", "path"})
)

// Metrics collects Prometheus metrics for each request.
func Metrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		httpRequestsInFlight.Inc()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := c.Response().StatusCode()
		path := c.Path()
		method := c.Method()

		httpRequestsTotal.WithLabelValues(method, path, statusCodeLabel(status)).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
		httpResponseSize.WithLabelValues(method, path).Observe(float64(len(c.Response().Body())))
		httpRequestsInFlight.Dec()

		return err
	}
}

func statusCodeLabel(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	default:
		return "5xx"
	}
}
