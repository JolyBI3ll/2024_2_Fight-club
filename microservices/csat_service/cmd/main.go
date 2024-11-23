package main

import (
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	grpcCsat "2024_2_FIGHT-CLUB/microservices/csat_service/controller"
	generatedCsat "2024_2_FIGHT-CLUB/microservices/csat_service/controller/gen"
	csatRepository "2024_2_FIGHT-CLUB/microservices/csat_service/repository"
	csatUseCase "2024_2_FIGHT-CLUB/microservices/csat_service/usecase"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	middleware.InitRedis()
	redisStore := session.NewRedisSessionStore(middleware.RedisClient)
	db := middleware.DbCSATConnect()

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

	csatRepo := csatRepository.NewCsatRepository(db)
	sessionService := session.NewSessionService(redisStore)
	csatUC := csatUseCase.NewCSATUseCase(csatRepo)
	csatServer := grpcCsat.NewGrpcCsatHandler(csatUC, sessionService, jwtToken)

	grpcServer := grpc.NewServer()
	generatedCsat.RegisterCsatServer(grpcServer, csatServer)

	listener, err := net.Listen("tcp", os.Getenv("CSAT_SERVICE_ADDRESS"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("CsatServer is listening on address: %s\n", os.Getenv("CSAT_SERVICE_ADDRESS"))
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
