package main

import (
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	grpcAd "2024_2_FIGHT-CLUB/microservices/ads_service/controller"
	generatedAds "2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	adRepository "2024_2_FIGHT-CLUB/microservices/ads_service/repository"
	adUseCase "2024_2_FIGHT-CLUB/microservices/ads_service/usecase"
	"context"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	middleware.InitRedis()
	redisStore := session.NewRedisSessionStore(middleware.RedisClient)
	db := middleware.DbConnect()
	minioService := middleware.MinioConnect()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Инициализация метрик
	metrics.InitMetrics()
	metrics.InitRepoMetric()
	// Экспозиция метрик на порту 9091
	go func() {
		http.Handle("/api/metrics", promhttp.Handler())
		log.Println("Metrics server is running on :9091")
		if err := http.ListenAndServe(":9091", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		if err := logger.SyncLoggers(); err != nil {
			log.Fatalf("Failed to sync loggers: %v", err)
		}
	}()

	jwtToken, err := middleware.NewJwtToken("secret-key")
	if err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}

	adsRepository := adRepository.NewAdRepository(db)
	sessionService := session.NewSessionService(redisStore)
	adsUseCase := adUseCase.NewAdUseCase(adsRepository, minioService)
	adsServer := grpcAd.NewGrpcAdHandler(sessionService, adsUseCase, jwtToken)
	adsUseCase.StartPriorityResetWorker(ctx)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.ChainUnaryInterceptors(
			middleware.RecoveryInterceptor,     // интерсептор для обработки паники
			middleware.UnaryMetricsInterceptor, // интерсептор для метрик
		)),
	)
	generatedAds.RegisterAdsServer(grpcServer, adsServer)

	listener, err := net.Listen("tcp", os.Getenv("ADS_SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("AdsServer is listening on address: %s\n", os.Getenv("ADS_SERVICE_ADDRESS"))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
