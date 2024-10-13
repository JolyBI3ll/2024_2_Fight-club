package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type AuthHandler struct {
	authUseCase    usecase.AuthUseCase
	sessionService *service.SessionService
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, sessionService *service.SessionService) *AuthHandler {
	return &AuthHandler{
		authUseCase:    authUseCase,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.authUseCase.RegisterUser(&user)
	if err != nil {
		if err.Error() == "username, password, and email are required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err.Error() == "user already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if err.Error() == "failed to generate error response" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	sessionID, err := h.sessionService.CreateSession(r, w, &user)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"user": map[string]string{
			"id":       user.UUID,
			"username": user.Username,
			"email":    user.Email,
		},
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) PutUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) GetUserById(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprint(w, `{"error": "Not Found"}`)
	if err != nil {
		return
	}
}

func SetUpRoutes(authHandler *AuthHandler) *mux.Router {
	router := mux.NewRouter()
	api := "/api"
	router.HandleFunc(api+"/auth/register", authHandler.RegisterUser).Methods("POST")
	router.HandleFunc(api+"/auth/login", authHandler.LoginUser).Methods("POST")
	router.HandleFunc(api+"/auth/logout", authHandler.LogoutUser).Methods("DELETE")
	router.HandleFunc(api+"/putUser", authHandler.PutUser).Methods("PUT")
	router.HandleFunc(api+"/getUserById", authHandler.GetUserById).Methods("GET")
	router.HandleFunc(api+"/getAllUser", authHandler.GetAllUser).Methods("GET")

	router.NotFoundHandler = authHandler

	return router
}
