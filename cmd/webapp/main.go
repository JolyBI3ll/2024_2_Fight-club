package main

import (
	generatedAds "2024_2_FIGHT-CLUB/ads_service/controller/grpc/gen"
	adHttpDelivery "2024_2_FIGHT-CLUB/ads_service/controller/http"
	adRepository "2024_2_FIGHT-CLUB/ads_service/repository"
	adUseCase "2024_2_FIGHT-CLUB/ads_service/usecase"
	generatedAuth "2024_2_FIGHT-CLUB/auth_service/controller/grpc/gen"
	authHttpDelivery "2024_2_FIGHT-CLUB/auth_service/controller/http"
	authRepository "2024_2_FIGHT-CLUB/auth_service/repository"
	authUseCase "2024_2_FIGHT-CLUB/auth_service/usecase"
	chatHttpDelivery "2024_2_FIGHT-CLUB/internal/chat/controller/http"
	chatRepository "2024_2_FIGHT-CLUB/internal/chat/repository"
	chatUseCase "2024_2_FIGHT-CLUB/internal/chat/usecase"
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

	adsConn, err := grpc.NewClient("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to AdsService: %v", err)
	}
	defer adsConn.Close()

	adsClient := generatedAds.NewAdsClient(adsConn)
	adsRepository := adRepository.NewAdRepository(db)
	adsUseCase := adUseCase.NewAdUseCase(adsRepository, minioService)
	adsHandler := adHttpDelivery.NewAdHandler(adsClient, adsUseCase, sessionService, jwtToken)

	citiesRepository := cityRepository.NewCityRepository(db)
	citiesUseCase := cityUseCase.NewCityUseCase(citiesRepository)
	cityHandler := cityHttpDelivery.NewCityHandler(citiesUseCase)

	chatsRepository := chatRepository.NewChatRepository(db)
	chatsUseCase := chatUseCase.NewChatService(chatsRepository)
	chatsHandler := chatHttpDelivery.NewChatController(chatsUseCase, sessionService)

	mainRouter := router.SetUpRoutes(authHandler, adsHandler, cityHandler, chatsHandler)
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
