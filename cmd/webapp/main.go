package main

import (
	authHttpDelivery "2024_2_FIGHT-CLUB/internal/auth/controller"
	authRepository "2024_2_FIGHT-CLUB/internal/auth/repository"
	authUseCase "2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service"
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
	auRepository := authRepository.NewAuthRepository(db)
	auUserCase := authUseCase.NewAuthUseCase(auRepository)
	sessionService := service.NewSessionService(store)
	authHandler := authHttpDelivery.NewAuthHandler(auUserCase, sessionService)
	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteStrictMode
	router := authHttpDelivery.SetUpRoutes(authHandler)
	http.Handle("/", enableCORS(router))
	fmt.Println("Starting server on port 8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}

//func main() {
//	_ = godotenv.Load()
//	store := sessions.NewCookieStore([]byte("super-secret-key"))
//	db := DbConnect()
//
//	auRepository := authRepository.NewAuthRepository(db)
//	auUserCase := authUseCase.NewAuthUseCase(auRepository)
//	sessionService := service.NewSessionService(store)
//	authHandler := authHttpDelivery.NewAuthHandler(auUserCase, sessionService)
//
//	store.Options.HttpOnly = true
//	store.Options.Secure = false
//	store.Options.SameSite = http.SameSiteStrictMode
//
//	router := authHttpDelivery.SetUpRoutes(authHandler)
//
//	router := mux.NewRouter()
//	authRepo := authRepository.NewAuthRepository(db)
//	authUCase := authUseCase.NewAuthUseCase(authRepo)
//	authHttpDelivery.NewAuthHandler(router, authUCase)
//	api := "/api"
//
//	router.HandleFunc(api+"/ads", controller.GetAllPlaces).Methods("GET")
//
//	router.HandleFunc(api+"/auth/register", controller.RegisterUser).Methods("POST")
//	router.HandleFunc(api+"/auth/login", controller.LoginUser).Methods("POST")
//	router.HandleFunc(api+"/auth/logout", controller.LogoutUser).Methods("DELETE")
//	router.HandleFunc(api+"/getAllUserData", controller.GetAllUserData).Methods("GET")
//	router.HandleFunc(api+"/getOneUserData", controller.GetOneUserData).Methods("GET")
//	router.HandleFunc(api+"/putUserData", controller.PutUserData).Methods("PUT")
//	router.HandleFunc(api+"/getSessionData", controller.GetSessionData).Methods("GET")
//
//	http.Handle("/", enableCORS(router))
//	fmt.Println("Starting server on port 8008")
//	if err := http.ListenAndServe(":8008", nil); err != nil {
//		fmt.Printf("Error on starting server: %s", err)
//	}
//}

func DbConnect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database")
	return db
}
