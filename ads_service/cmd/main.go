package main

import (
	grpcAd "2024_2_FIGHT-CLUB/ads_service/controller/grpc"
	generatedAds "2024_2_FIGHT-CLUB/ads_service/controller/grpc/gen"
	adRepository "2024_2_FIGHT-CLUB/ads_service/repository"
	adUseCase "2024_2_FIGHT-CLUB/ads_service/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	middleware.InitRedis()
	redisStore := session.NewRedisSessionStore(middleware.RedisClient)
	db := middleware.DbConnect()
	minioService := middleware.MinioConnect()

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

	grpcServer := grpc.NewServer()
	generatedAds.RegisterAdsServer(grpcServer, adsServer)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("AdsServer is listening on port 50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
