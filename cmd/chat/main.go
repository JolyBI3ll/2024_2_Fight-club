package main

import (
	grpcAuth "2024_2_FIGHT-CLUB/internal/auth/controller/grpc"
	generatedAuth "2024_2_FIGHT-CLUB/internal/auth/controller/grpc/gen"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
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

	jwtToken, err := middleware.NewJwtToken("secret-key")
	if err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}

	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		if err := logger.SyncLoggers(); err != nil {
			log.Fatalf("Failed to sync loggers: %v", err)
		}
	}()

	authServer := grpcAuth.NewGrpcAuthHandler(auUseCase, sessionService, jwtToken)

	grpcServer := grpc.NewServer()
	generatedAuth.RegisterAuthServer(grpcServer, authServer)

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("AuthService is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
