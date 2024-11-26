package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GrpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Histogram of response duration for gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	GrpcErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_service_errors_total",
			Help: "Total number of errors in the grpc service",
		},
		[]string{"method", "status", "error"},
	)
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "url", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "url"},
	)

	HttpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_service_errors_total",
			Help: "Total number of errors in the htpp service",
		},
		[]string{"method", "url", "status", "error"},
	)
)

// InitMetrics регистрирует метрики
func InitMetrics() {
	prometheus.MustRegister(GrpcRequestsTotal)
	prometheus.MustRegister(GrpcRequestDuration)
	prometheus.MustRegister(GrpcErrorsTotal)
}

func InitHttpMetric() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpErrorsTotal)
}
