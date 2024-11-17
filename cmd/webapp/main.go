package main

import (
	adHttpDelivery "2024_2_FIGHT-CLUB/internal/ads/controller"
	adRepository "2024_2_FIGHT-CLUB/internal/ads/repository"
	adUseCase "2024_2_FIGHT-CLUB/internal/ads/usecase"
	generatedAuth "2024_2_FIGHT-CLUB/internal/auth/controller/grpc/gen"
	authHttpDelivery "2024_2_FIGHT-CLUB/internal/auth/controller/http"
	authRepository "2024_2_FIGHT-CLUB/internal/auth/repository"
	authUseCase "2024_2_FIGHT-CLUB/internal/auth/usecase"
	cityHttpDelivery "2024_2_FIGHT-CLUB/internal/cities/controller"
	cityRepository "2024_2_FIGHT-CLUB/internal/cities/repository"
	cityUseCase "2024_2_FIGHT-CLUB/internal/cities/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/router"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()
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
		err := logger.SyncLoggers()
		if err != nil {
			log.Fatalf("Failed to sync loggers: %v", err)
		}
	}()

	authConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) // Укажите адрес AuthService
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	sessionService := session.NewSessionService(redisStore)

	authClient := generatedAuth.NewAuthClient(authConn)
	auRepository := authRepository.NewAuthRepository(db)
	auUseCase := authUseCase.NewAuthUseCase(auRepository, minioService)
	authHandler := authHttpDelivery.NewAuthHandler(authClient, auUseCase, sessionService, jwtToken)

	adsRepository := adRepository.NewAdRepository(db)
	adsUseCase := adUseCase.NewAdUseCase(adsRepository, minioService)
	adsHandler := adHttpDelivery.NewAdHandler(adsUseCase, sessionService, jwtToken)

	citiesRepository := cityRepository.NewCityRepository(db)
	citiesUseCase := cityUseCase.NewCityUseCase(citiesRepository)
	cityHandler := cityHttpDelivery.NewCityHandler(citiesUseCase)

	mainRouter := router.SetUpRoutes(authHandler, adsHandler, cityHandler)
	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.RateLimitMiddleware)
	http.Handle("/", middleware.EnableCORS(mainRouter))
	if os.Getenv("HTTPS") == "TRUE" {
		fmt.Println("Starting HTTPS server on port 8008")
		if err := http.ListenAndServeTLS("0.0.0.0:8008", "ssl/pootnick.crt", "ssl/pootnick.key", nil); err != nil {
			fmt.Printf("Error on starting server: %s", err)
		}
	} else {
		fmt.Println("Starting HTTP server on port 8008")
		if err := http.ListenAndServe("0.0.0.0:8008", nil); err != nil {
			fmt.Printf("Error on starting server: %s", err)
		}
	}

}
