package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// contains all Prometheus metrics for the service.
type Metrics struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
	ReviewsCreatedTotal *prometheus.CounterVec
}

// creates and registers the Prometheus metrics.
func NewMetrics(reg prometheus.Registerer) *Metrics {
	return &Metrics{
		HttpRequestsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "code", "path"},
		),
		HttpRequestDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests.",
				Buckets: prometheus.DefBuckets, // Default buckets: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
			},
			[]string{"method", "path"},
		),
		ReviewsCreatedTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "reviews_created_total",
				Help: "Total number of created reviews by sentiment.",
			},
			[]string{"sentiment"},
		),
	}
}
