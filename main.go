package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	http.Handle("/", enableCORS(router))
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error on starting server: %s", err)
	}
}
