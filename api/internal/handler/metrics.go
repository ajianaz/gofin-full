package handler

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

// MetricsHandler serves Prometheus metrics.
type MetricsHandler struct{}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// Prometheus exposes Prometheus metrics at /metrics.
func (h *MetricsHandler) Prometheus(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.Set("X-Content-Type-Options", "nosniff")

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return c.Status(500).SendString("failed to gather metrics")
	}

	var buf bytes.Buffer
	for _, mf := range mfs {
		_, _ = expfmt.MetricFamilyToText(&buf, mf)
	}

	return c.Send(buf.Bytes())
}
