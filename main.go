package main

import (
	"2024_2_FIGHT-CLUB/dsn"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
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

var db *gorm.DB

func main() {
	_ = godotenv.Load()
	var err error
	db, err = gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to database")

	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteStrictMode

	router := mux.NewRouter()
	api := "/api"

	router.HandleFunc(api+"/ads", getAllPlaces).Methods("GET")

	router.HandleFunc(api+"/auth/register", registerUser).Methods("POST")
	router.HandleFunc(api+"/auth/login", loginUser).Methods("POST")
	router.HandleFunc(api+"/auth/logout", logoutUser).Methods("DELETE")

	router.HandleFunc(api+"/getSessionData", getSessionData).Methods("GET")
	router.HandleFunc(api+"/getAllUserData", GetAllUserData).Methods("GET")
	router.HandleFunc(api+"/getOneUserData", GetOneUserData).Methods("GET")
	router.HandleFunc(api+"/putUSerData", PutUserData).Methods("GET")

	http.Handle("/", enableCORS(router))
	fmt.Println("Starting server on port 8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
