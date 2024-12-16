package main

import (
	adHttpDelivery "2024_2_FIGHT-CLUB/internal/ads/controller"
	authHttpDelivery "2024_2_FIGHT-CLUB/internal/auth/controller"
	chatHttpDelivery "2024_2_FIGHT-CLUB/internal/chat/controller"
	chatRepository "2024_2_FIGHT-CLUB/internal/chat/repository"
	chatUseCase "2024_2_FIGHT-CLUB/internal/chat/usecase"
	cityHttpDelivery "2024_2_FIGHT-CLUB/internal/cities/controller"
	reviewContoller "2024_2_FIGHT-CLUB/internal/reviews/contoller"
	reviewRepository "2024_2_FIGHT-CLUB/internal/reviews/repository"
	reviewUsecase "2024_2_FIGHT-CLUB/internal/reviews/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/router"
	"2024_2_FIGHT-CLUB/internal/service/session"
	generatedAds "2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	generatedAuth "2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	generatedCity "2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	metrics.InitHttpMetric()
	metrics.InitRepoMetric()

	authAdress := os.Getenv("AUTH_SERVICE_ADDRESS")
	if authAdress == "" {
		log.Fatalf("AUTH_SERVICE_ADDRESS is not set")
	}
	authConn, err := grpc.NewClient(authAdress, grpc.WithTransportCredentials(insecure.NewCredentials())) // Укажите адрес AuthService
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
	}
	defer authConn.Close()

	adsAdress := os.Getenv("ADS_SERVICE_ADDRESS")
	if adsAdress == "" {
		log.Fatalf("ADS_SERVICE_ADDRESS is not set")
	}
	adsConn, err := grpc.NewClient(adsAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to AdsService: %v", err)
	}
	defer adsConn.Close()

	cityAdress := os.Getenv("CITY_SERVICE_ADDRESS")
	if cityAdress == "" {
		log.Fatalf("CITY_SERVICE_ADDRESS is not set")
	}
	cityConn, err := grpc.NewClient(cityAdress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to AdsService: %v", err)
	}
	defer cityConn.Close()

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

	reviewsRepository := reviewRepository.NewReviewRepository(db)
	reviewsUsecase := reviewUsecase.NewReviewUsecase(reviewsRepository)
	reviewsHandler := reviewContoller.NewReviewHandler(reviewsUsecase, sessionService, jwtToken)

	mainRouter := router.SetUpRoutes(authHandler, adsHandler, cityHandler, chatsHandler, reviewsHandler)
	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.RateLimitMiddleware)
	http.Handle("/", middleware.RecoverWrap(middleware.EnableCORS(mainRouter)))
	if os.Getenv("HTTPS") == "TRUE" {
		fmt.Printf("Starting HTTPS server on address %s\n", os.Getenv("BACKEND_URL"))
		if err := http.ListenAndServeTLS(os.Getenv("BACKEND_URL"), "ssl/pootnick.crt", "ssl/pootnick.key", nil); err != nil {
			fmt.Printf("Error on starting server: %s", err)
		}
	} else {
		fmt.Printf("Starting HTTP server on adress %s\n", os.Getenv("BACKEND_URL"))
		if err := http.ListenAndServe(os.Getenv("BACKEND_URL"), nil); err != nil {
			fmt.Printf("Error on starting server: %s", err)
		}
	}

}
