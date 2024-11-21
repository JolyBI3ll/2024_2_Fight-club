package main

import (
	adHttpDelivery "2024_2_FIGHT-CLUB/internal/ads/controller"
	authHttpDelivery "2024_2_FIGHT-CLUB/internal/auth/controller"
	chatHttpDelivery "2024_2_FIGHT-CLUB/internal/chat/controller"
	chatRepository "2024_2_FIGHT-CLUB/internal/chat/repository"
	chatUseCase "2024_2_FIGHT-CLUB/internal/chat/usecase"
	cityHttpDelivery "2024_2_FIGHT-CLUB/internal/cities/controller"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/router"
	"2024_2_FIGHT-CLUB/internal/service/session"
	generatedAds "2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	generatedAuth "2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	generatedCity "2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
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

	adsConn, err := grpc.NewClient("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to AdsService: %v", err)
	}
	defer adsConn.Close()

	cityConn, err := grpc.NewClient("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to AdsService: %v", err)
	}
	defer adsConn.Close()

	sessionService := session.NewSessionService(redisStore)

	authClient := generatedAuth.NewAuthClient(authConn)
	authHandler := authHttpDelivery.NewAuthHandler(authClient, sessionService, jwtToken)

	adsClient := generatedAds.NewAdsClient(adsConn)
	adsHandler := adHttpDelivery.NewAdHandler(adsClient, sessionService, jwtToken)

	cityClient := generatedCity.NewCityServiceClient(cityConn)
	cityHandler := cityHttpDelivery.NewCityHandler(cityClient)

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
