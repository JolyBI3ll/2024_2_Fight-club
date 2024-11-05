package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/internal/service/validation"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"regexp"
	"time"
)

type AuthHandler struct {
	authUseCase    usecase.AuthUseCase
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *AuthHandler {
	return &AuthHandler{
		authUseCase:    authUseCase,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

const requestTimeout = 5 * time.Second

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, requestTimeout)
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sanitizer := bluemonday.UGCPolicy()
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

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
		h.handleError(w, err, requestID)
		return
	}

	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s]*$`)
	if !validCharPattern.MatchString(creds.Avatar) ||
		!validCharPattern.MatchString(creds.UUID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(creds.Username) > maxLen || len(creds.Email) > maxLen || len(creds.Password) > maxLen || len(creds.Name) > maxLen || len(creds.Avatar) > maxLen || len(creds.UUID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	creds.Avatar = sanitizer.Sanitize(creds.Avatar)
	creds.Username = sanitizer.Sanitize(creds.Username)
	creds.Email = sanitizer.Sanitize(creds.Email)
	creds.Password = sanitizer.Sanitize(creds.Password)
	creds.UUID = sanitizer.Sanitize(creds.UUID)
	creds.Name = sanitizer.Sanitize(creds.Name)

	err := h.authUseCase.RegisterUser(ctx, &creds)

	if err != nil {
		logger.AccessLogger.Error("Failed to register user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	userSession, err := h.sessionService.CreateSession(ctx, r, w, &creds)
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
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
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
		h.handleError(w, err, requestID)
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
	sanitizer := bluemonday.UGCPolicy()
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

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
		h.handleError(w, err, requestID)
		return
	}

	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s]*$`)
	if !validCharPattern.MatchString(creds.Email) ||
		!validCharPattern.MatchString(creds.Name) ||
		!validCharPattern.MatchString(creds.Avatar) ||
		!validCharPattern.MatchString(creds.UUID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(creds.Username) > maxLen || len(creds.Email) > maxLen || len(creds.Password) > maxLen || len(creds.Name) > maxLen || len(creds.Avatar) > maxLen || len(creds.UUID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	creds.Avatar = sanitizer.Sanitize(creds.Avatar)
	creds.Username = sanitizer.Sanitize(creds.Username)
	creds.Email = sanitizer.Sanitize(creds.Email)
	creds.Password = sanitizer.Sanitize(creds.Password)
	creds.UUID = sanitizer.Sanitize(creds.UUID)
	creds.Name = sanitizer.Sanitize(creds.Name)

	requestedUser, err := h.authUseCase.LoginUser(ctx, &creds)

	if err != nil {
		logger.AccessLogger.Error("Failed to login user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	userSession, err := h.sessionService.CreateSession(ctx, r, w, requestedUser)

	if err != nil {
		logger.AccessLogger.Error("Failed create session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	csrfToken, _ := r.Cookie("csrf_token")
	if csrfToken != nil {
		logger.AccessLogger.Error("csrf_token already exists",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, errors.New("csrf_token already exists"), requestID)
		return
	}

	tokenExpTime := time.Now().Add(24 * time.Hour).Unix()
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
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
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

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received LogoutUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		h.handleError(w, errors.New("Missing X-CSRF-Token header"), requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Invalid JWT token"), requestID)
		return
	}

	if err := h.sessionService.LogoutSession(ctx, r, w); err != nil {

		logger.AccessLogger.Error("Failed to logout user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logoutResponse := map[string]string{"response": "Logout successfully"}
	if err := json.NewEncoder(w).Encode(logoutResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode logout response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
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
	sanitizer := bluemonday.UGCPolicy()
	requestID := middleware.GetRequestID(r.Context())

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received PutUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	authHeader := r.Header.Get("X-CSRF-Token")
	if authHeader == "" {
		logger.AccessLogger.Warn("Failed to X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("Missing X-CSRF-Token header")),
		)
		h.handleError(w, errors.New("Missing X-CSRF-Token header"), requestID)
		return
	}

	tokenString := authHeader[len("Bearer "):]
	_, err := h.jwtToken.Validate(tokenString)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, errors.New("Invalid JWT token"), requestID)
		return
	}

	r.ParseMultipartForm(10 << 20)

	metadata := r.FormValue("metadata")

	var creds domain.User
	if err := json.Unmarshal([]byte(metadata), &creds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.AccessLogger.Warn("Failed to parse metadata",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, errors.New("Invalid metadata JSON"), requestID)
		return
	}

	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s]*$`)
	if !validCharPattern.MatchString(creds.Username) ||
		!validCharPattern.MatchString(creds.Email) ||
		!validCharPattern.MatchString(creds.Password) ||
		!validCharPattern.MatchString(creds.Name) ||
		!validCharPattern.MatchString(creds.Avatar) ||
		!validCharPattern.MatchString(creds.UUID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(creds.Username) > maxLen || len(creds.Email) > maxLen || len(creds.Password) > maxLen || len(creds.Name) > maxLen || len(creds.Avatar) > maxLen || len(creds.UUID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	creds.Avatar = sanitizer.Sanitize(creds.Avatar)
	creds.Username = sanitizer.Sanitize(creds.Username)
	creds.Email = sanitizer.Sanitize(creds.Email)
	creds.Password = sanitizer.Sanitize(creds.Password)
	creds.UUID = sanitizer.Sanitize(creds.UUID)
	creds.Name = sanitizer.Sanitize(creds.Name)

	var avatar *multipart.FileHeader
	if len(r.MultipartForm.File["avatar"]) > 0 {
		avatar = r.MultipartForm.File["avatar"][0]

		if err := validation.ValidateImage(avatar, 5<<20, []string{"image/jpeg", "image/png", "image/jpg"}, 2000, 2000); err != nil {
			logger.AccessLogger.Warn("Invalid size, type or resolution of image", zap.String("request_id", requestID), zap.Error(err))
			h.handleError(w, errors.New("Invalid size, type or resolution of image"), requestID)
			return
		}
	}

	userID, err := h.sessionService.GetUserID(ctx, r)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user ID from session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	if err := h.authUseCase.PutUser(ctx, &creds, userID, avatar); err != nil {
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
		h.handleError(w, err, requestID)
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
	sanitizer := bluemonday.UGCPolicy()

	userId := mux.Vars(r)["userId"]
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(userId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input contains invalid characters"), requestID)
		return
	}

	if len(userId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		h.handleError(w, errors.New("Input exceeds character limit"), requestID)
		return
	}

	userId = sanitizer.Sanitize(userId)

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetUserById request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	w.Header().Set("Content-Type", "application/json")
	user, err := h.authUseCase.GetUserById(ctx, userId)
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
		h.handleError(w, err, requestID)
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

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetAllUsers request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	w.Header().Set("Content-Type", "application/json")
	users, err := h.authUseCase.GetAllUser(ctx)
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
		h.handleError(w, err, requestID)
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

	ctx, cancel := withTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetSessionData request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionData, err := h.sessionService.GetSessionData(ctx, r)

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
		h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetSessionData request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) RefreshCsrfToken(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received RefreshCsrfToken request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	userSession, err := h.sessionService.GetSession(r.Context(), r)
	if err != nil {
		logger.AccessLogger.Error("Failed to session",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	newCsrfToken, err := h.jwtToken.Create(userSession, time.Now().Add(1*time.Hour).Unix())
	if err != nil {
		logger.AccessLogger.Error("Failed to generate CSRF",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		http.Error(w, "Failed to create CSRF token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    newCsrfToken,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"csrf_token": newCsrfToken})
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed RefreshCsrfToken request",
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
		"invalid credentials", "csrf_token already exists", "Input contains invalid characters",
		"Input exceeds character limit", "Invalid size, type or resolution of image":
		w.WriteHeader(http.StatusBadRequest)
	case "user already exists",
		"session already exists":
		w.WriteHeader(http.StatusConflict)
	case "no active session", "already logged in":
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
