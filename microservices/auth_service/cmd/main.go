package main

import (
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	grpcAuth "2024_2_FIGHT-CLUB/microservices/auth_service/controller"
	generatedAuth "2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	authRepository "2024_2_FIGHT-CLUB/microservices/auth_service/repository"
	authUseCase "2024_2_FIGHT-CLUB/microservices/auth_service/usecase"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализация зависимостей
	middleware.InitRedis()
	redisStore := session.NewRedisSessionStore(middleware.RedisClient)
	db := middleware.DbConnect()
	minioService := middleware.MinioConnect()

	// Инициализация метрик
	metrics.InitMetrics()

	// Экспозиция метрик на порту 9092
	go func() {
		http.Handle("/api/metrics", promhttp.Handler())
		log.Println("Metrics server is running on :9092")
		if err := http.ListenAndServe(":9092", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Создание JWT сервиса
	jwtToken, err := middleware.NewJwtToken("secret-key")
	if err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}

	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		if err := logger.SyncLoggers(); err != nil {
			log.Fatalf("Failed to sync loggers: %v", err)
		}
	}()

	sessionService := session.NewSessionService(redisStore)
	auRepository := authRepository.NewAuthRepository(db)
	auUseCase := authUseCase.NewAuthUseCase(auRepository, minioService)
	authServer := grpcAuth.NewGrpcAuthHandler(auUseCase, sessionService, jwtToken)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryMetricsInterceptor),
	)
	generatedAuth.RegisterAuthServer(grpcServer, authServer)

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", os.Getenv("AUTH_SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("Failed to listen on address: %s %v", os.Getenv("AUTH_SERVICE_ADDRESS"), err)
	}

	log.Printf("AuthService is running on address: %s\n", os.Getenv("AUTH_SERVICE_ADDRESS"))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
