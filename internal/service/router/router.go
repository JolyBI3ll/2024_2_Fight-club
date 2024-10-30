package router

import (
	ads "2024_2_FIGHT-CLUB/internal/ads/controller"
	auth "2024_2_FIGHT-CLUB/internal/auth/controller"
	"github.com/gorilla/mux"
)

func SetUpRoutes(authHandler *auth.AuthHandler, adsHandler *ads.AdHandler) *mux.Router {
	router := mux.NewRouter()
	api := "/api"

	// User Authentication Routes
	router.HandleFunc(api+"/auth/register", authHandler.RegisterUser).Methods("POST") // Register a new user
	router.HandleFunc(api+"/auth/login", authHandler.LoginUser).Methods("POST")       // Login user
	router.HandleFunc(api+"/auth/logout", authHandler.LogoutUser).Methods("DELETE")   // Logout user

	// User Management Routes
	router.HandleFunc(api+"/users/{userId}", authHandler.PutUser).Methods("PUT")     // Update user
	router.HandleFunc(api+"/users/{userId}", authHandler.GetUserById).Methods("GET") // Get user by ID
	router.HandleFunc(api+"/users", authHandler.GetAllUsers).Methods("GET")          // Get all users
	router.HandleFunc(api+"/session", authHandler.GetSessionData).Methods("GET")     // Get session data

	// Ad Management Routes
	router.HandleFunc(api+"/ads", adsHandler.GetAllPlaces).Methods("GET")                   // Get all ads
	router.HandleFunc(api+"/ads/{adId}", adsHandler.GetOnePlace).Methods("GET")             // Get ad by ID
	router.HandleFunc(api+"/ads", adsHandler.CreatePlace).Methods("POST")                   // Create a new ad
	router.HandleFunc(api+"/ads/{adId}", adsHandler.UpdatePlace).Methods("PUT")             // Update ad by ID
	router.HandleFunc(api+"/ads/{adId}", adsHandler.DeletePlace).Methods("DELETE")          // Delete ad by ID
	router.HandleFunc(api+"/ads/cities/{city}", adsHandler.GetPlacesPerCity).Methods("GET") // Get ads by city

	// CSRF Token Route
	router.HandleFunc(api+"/csrf/refresh", authHandler.RefreshCsrfToken).Methods("GET") // Refresh CSRF token

	return router
}
