package router

import (
	ads "2024_2_FIGHT-CLUB/internal/ads/controller"
	auth "2024_2_FIGHT-CLUB/internal/auth/controller"
	"github.com/gorilla/mux"
)

func SetUpRoutes(authHandler *auth.AuthHandler, adsHandler *ads.AdHandler) *mux.Router {
	router := mux.NewRouter()
	api := "/api"

	router.HandleFunc(api+"/auth/register", authHandler.RegisterUser).Methods("POST")
	router.HandleFunc(api+"/auth/login", authHandler.LoginUser).Methods("POST")

	router.HandleFunc(api+"/auth/logout", authHandler.LogoutUser).Methods("DELETE")

	router.HandleFunc(api+"/putUser", authHandler.PutUser).Methods("PUT")

	router.HandleFunc(api+"/getUserById", authHandler.GetUserById).Methods("GET")
	router.HandleFunc(api+"/getAllUsers", authHandler.GetAllUsers).Methods("GET")
	router.HandleFunc(api+"/getSessionData", authHandler.GetSessionData).Methods("GET")

	router.HandleFunc(api+"/ads", adsHandler.GetAllPlaces).Methods("GET")
	router.HandleFunc(api+"/ads/{adId}", adsHandler.GetOnePlace).Methods("GET")
	router.HandleFunc(api+"/createAd", adsHandler.CreatePlace).Methods("POST")
	router.HandleFunc(api+"/updateAd/{adId}", adsHandler.UpdatePlace).Methods("PUT")
	router.HandleFunc(api+"/deleteAd/{adId}", adsHandler.DeletePlace).Methods("DELETE")
	return router
}
