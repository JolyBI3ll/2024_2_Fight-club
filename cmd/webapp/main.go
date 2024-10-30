package main

import (
	adHttpDelivery "2024_2_FIGHT-CLUB/internal/ads/controller"
	adRepository "2024_2_FIGHT-CLUB/internal/ads/repository"
	adUseCase "2024_2_FIGHT-CLUB/internal/ads/usecase"
	authHttpDelivery "2024_2_FIGHT-CLUB/internal/auth/controller"
	authRepository "2024_2_FIGHT-CLUB/internal/auth/repository"
	authUseCase "2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/router"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	_ = godotenv.Load()
	store := sessions.NewCookieStore([]byte("super-secret-key"))
	db := middleware.DbConnect()
	minioService := middleware.MinioConnect()
	jwtToken, err := middleware.NewJwtToken("secret-key")
	if err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}

	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	sessionService := session.NewSessionService(store)

	auRepository := authRepository.NewAuthRepository(db)
	auUseCase := authUseCase.NewAuthUseCase(auRepository, minioService)
	authHandler := authHttpDelivery.NewAuthHandler(auUseCase, sessionService, jwtToken)

	adsRepository := adRepository.NewAdRepository(db)
	adsUseCase := adUseCase.NewAdUseCase(adsRepository, minioService)
	adsHandler := adHttpDelivery.NewAdHandler(adsUseCase, sessionService, jwtToken)

	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteStrictMode

	mainRouter := router.SetUpRoutes(authHandler, adsHandler)
	mainRouter.Use(middleware.RequestIDMiddleware)
	mainRouter.Use(middleware.RateLimitMiddleware)
	http.Handle("/", middleware.EnableCORS(mainRouter))
	fmt.Println("Starting server on port 8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
