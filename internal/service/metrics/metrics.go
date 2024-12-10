package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"regexp"
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
		[]string{"method", "url", "status", "remote_ip"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "url", "remote_ip"},
	)

	HttpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_service_errors_total",
			Help: "Total number of errors in the htpp service",
		},
		[]string{"method", "url", "status", "error", "remote_ip"},
	)
)

var (
	RepoRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repository_sql_total",
			Help: "Total number of sql request",
		},
		[]string{"method", "status"},
	)
	RepoRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "repo_request_duration_second",
			Help:    "Histogram of request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"})
	RepoErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "repo_request_errors_total",
			Help: "Total number of errors in the repo",
		},
		[]string{"method", "status", "error"},
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

func InitRepoMetric() {
	prometheus.MustRegister(RepoRequestTotal)
	prometheus.MustRegister(RepoRequestDuration)
	prometheus.MustRegister(RepoErrorsTotal)
}

func SanitizeUserIdPath(path string) string {
	re := regexp.MustCompile(`[0-9a-fA-F-]{36}`)
	return re.ReplaceAllString(path, "{userId}")
}

func SanitizeAdIdPath(path string) string {
	re := regexp.MustCompile(`[0-9a-fA-F-]{36}`)
	return re.ReplaceAllString(path, "{adId}")
}
