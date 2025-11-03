package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler exposes the Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
