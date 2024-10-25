package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"encoding/json"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"time"
)

type AuthHandler struct {
	authUseCase    usecase.AuthUseCase
	sessionService *session.ServiceSession
	jwtToken       *middleware.JwtToken
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, sessionService *session.ServiceSession, jwtToken *middleware.JwtToken) *AuthHandler {
	return &AuthHandler{
		authUseCase:    authUseCase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received RegisterUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var creds domain.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		logger.AccessLogger.Error("Failed to decode request body",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.authUseCase.RegisterUser(r.Context(), &creds)
	if err != nil {
		logger.AccessLogger.Error("Failed to register user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	userSession, err := h.sessionService.CreateSession(r.Context(), r, w, &creds)
	if err != nil {
		logger.AccessLogger.Error("Failed create session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	tokenExpTime := time.Now().Add(24 * time.Hour).Unix() // например, срок действия 24 часа
	jwtToken, err := h.jwtToken.Create(userSession, tokenExpTime)
	if err != nil {
		logger.AccessLogger.Error("Failed to create JWT token",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
	})

	response := map[string]interface{}{
		"session_id": userSession.Values["session_id"],
		"user": map[string]string{
			"id":       creds.UUID,
			"username": creds.Username,
			"email":    creds.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed RegisterUser request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusCreated),
	)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received LoginUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var creds domain.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		logger.AccessLogger.Error("Failed to decode request body",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestedUser, err := h.authUseCase.LoginUser(r.Context(), &creds)
	if err != nil {
		logger.AccessLogger.Error("Failed to login user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	userSession, err := h.sessionService.CreateSession(r.Context(), r, w, requestedUser)
	if err != nil {
		logger.AccessLogger.Error("Failed create session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	tokenExpTime := time.Now().Add(24 * time.Hour).Unix() // срок действия 24 часа
	jwtToken, err := h.jwtToken.Create(userSession, tokenExpTime)
	if err != nil {
		logger.AccessLogger.Error("Failed to create JWT token",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    jwtToken,
		HttpOnly: true,
		Secure:   true,
	})

	response := map[string]interface{}{
		"session_id": userSession.Values["session_id"],
		"user": map[string]interface{}{
			"id":       requestedUser.UUID,
			"username": requestedUser.Username,
			"email":    requestedUser.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed LoginUser request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received LogoutUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	//authHeader := r.Header.Get("X-CSRF-Token")
	//if authHeader == "" {
	//	http.Error(w, "Missing X-CSRF-Token header", http.StatusUnauthorized)
	//	return
	//}
	//
	//tokenString := authHeader[len("Bearer "):]
	//_, err := h.jwtToken.Validate(tokenString)
	//if err != nil {
	//	logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
	//	http.Error(w, "Invalid token", http.StatusUnauthorized)
	//	return
	//}

	if err := h.sessionService.LogoutSession(r.Context(), r, w); err != nil {
		logger.AccessLogger.Error("Failed to logout user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logoutResponse := map[string]string{"response": "Logout successfully"}
	if err := json.NewEncoder(w).Encode(logoutResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode logout response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed LogoutUser request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) PutUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received PutUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	r.ParseMultipartForm(10 << 20)

	metadata := r.FormValue("metadata")

	var creds domain.User
	if err := json.Unmarshal([]byte(metadata), &creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Invalid metadata JSON", http.StatusBadRequest)
		return
	}
	var avatar *multipart.FileHeader
	if len(r.MultipartForm.File["avatar"]) > 0 {
		avatar = r.MultipartForm.File["avatar"][0]
	}

	userID, err := h.sessionService.GetUserID(r.Context(), r, w)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID from session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}
	if err := h.authUseCase.PutUser(r.Context(), &creds, userID, avatar); err != nil {
		logger.AccessLogger.Error("Failed to update user data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("Update successful"); err != nil {
		logger.AccessLogger.Error("Failed to encode update response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed PutUser request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received GetUserById request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	w.Header().Set("Content-Type", "application/json")
	userID, err := h.sessionService.GetUserID(r.Context(), r, w)
	if err != nil {
		logger.AccessLogger.Error("Failed to get user ID from session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	user, err := h.authUseCase.GetUserById(r.Context(), userID)
	if err != nil {
		logger.AccessLogger.Error("Failed to get user by id",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		logger.AccessLogger.Error("Failed to encode user data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetUserById request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received GetAllUsers request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	w.Header().Set("Content-Type", "application/json")
	users, err := h.authUseCase.GetAllUser(r.Context())
	if err != nil {
		logger.AccessLogger.Error("Failed to get all users data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	response := map[string]interface{}{
		"users": users,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode users response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetAllUsers request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) GetSessionData(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received GetSessionData request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionData, err := h.sessionService.GetSessionData(r.Context(), r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sessionData); err != nil {
		logger.AccessLogger.Error("Failed to encode session data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetSessionData request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

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
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}
