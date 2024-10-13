package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authUseCase    usecase.AuthUseCase
	sessionService *session.ServiceSession
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, sessionService *session.ServiceSession) *AuthHandler {
	return &AuthHandler{
		authUseCase:    authUseCase,
		sessionService: sessionService,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var creds domain.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.authUseCase.RegisterUser(&creds)
	if err != nil {
		h.handleError(w, err)
		return
	}

	sessionID, err := h.sessionService.CreateSession(r, w, &creds)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"user": map[string]string{
			"id":       creds.UUID,
			"username": creds.Username,
			"email":    creds.Email,
		},
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds domain.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestedUser, err := h.authUseCase.LoginUser(&creds)
	if err != nil {
		h.handleError(w, err)
		return
	}

	sessionID, err := h.sessionService.CreateSession(r, w, requestedUser)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{
		"session_id": sessionID,
		"user": map[string]interface{}{
			"id":       requestedUser.UUID,
			"username": requestedUser.Username,
			"email":    requestedUser.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	if err := h.sessionService.LogoutSession(r, w); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logoutResponse := map[string]string{"response": "Logout successfully"}
	if err := json.NewEncoder(w).Encode(logoutResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) PutUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.sessionService.GetUserID(r, w)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err := h.authUseCase.PutUser(&user, userID); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("Update successful"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := h.sessionService.GetUserID(r, w)
	if err != nil {
		h.handleError(w, err)
		return
	}
	user, err := h.authUseCase.GetUserById(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.authUseCase.GetAllUser()
	if err != nil {
		h.handleError(w, err)
		return
	}
	response := map[string]interface{}{
		"users": users,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) GetSessionData(w http.ResponseWriter, r *http.Request) {
	sessionData, err := h.sessionService.GetSessionData(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sessionData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}

	switch err.Error() {
	case "username, password, and email are required",
		"username and password are required",
		"invalid credentials":
		w.WriteHeader(http.StatusBadRequest)
	case "user already exists",
		"session already exists":
		w.WriteHeader(http.StatusConflict)
	case "no active session":
		w.WriteHeader(http.StatusUnauthorized)
	case "user not found":
		w.WriteHeader(http.StatusNotFound)
	case "failed to generate error response",
		"there is none user in db",
		"failed to generate session id",
		"failed to save sessions":
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}
