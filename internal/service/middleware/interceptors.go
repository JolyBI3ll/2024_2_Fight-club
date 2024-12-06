package middleware

import (
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func UnaryMetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	statusCode := "success"
	if err != nil {
		statusCode = "error"
	}
	metrics.GrpcRequestsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
	metrics.GrpcRequestDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

	// Метрики ошибок
	if err != nil {
		grpcStatus, _ := status.FromError(err) // gRPC статус ошибки
		errorMsg := grpcStatus.Message()
		metrics.GrpcErrorsTotal.WithLabelValues(info.FullMethod, grpcStatus.Code().String(), errorMsg).Inc()
	}

	return resp, err
}
