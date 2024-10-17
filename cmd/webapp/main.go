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
	"2024_2_FIGHT-CLUB/module/dsn"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Set-Cookie")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()
	store := sessions.NewCookieStore([]byte("super-secret-key"))
	db := DbConnect()
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	sessionService := session.NewSessionService(store)

	auRepository := authRepository.NewAuthRepository(db)
	auUseCase := authUseCase.NewAuthUseCase(auRepository)
	authHandler := authHttpDelivery.NewAuthHandler(auUseCase, sessionService)

	adsRepository := adRepository.NewAdRepository(db)
	adsUseCase := adUseCase.NewAdUseCase(adsRepository)
	adsHandler := adHttpDelivery.NewAdHandler(adsUseCase, sessionService)

	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteStrictMode

	mainRouter := router.SetUpRoutes(authHandler, adsHandler)
	mainRouter.Use(middleware.RequestIDMiddleware)
	http.Handle("/", enableCORS(mainRouter))
	fmt.Println("Starting server on port 8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}

func DbConnect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database")
	return db
}
