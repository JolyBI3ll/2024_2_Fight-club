package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/internal/service/utils"
	"2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	client         gen.AuthClient
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewAuthHandler(client gen.AuthClient, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *AuthHandler {
	return &AuthHandler{
		client:         client,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusCreated
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusCreated {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received RegisterUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var creds domain.User
	if err = easyjson.UnmarshalFromReader(r.Body, &creds); err != nil {
		logger.AccessLogger.Error("Failed to decode request body",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.RegisterUser(ctx, &gen.RegisterUserRequest{
		Username: creds.Username,
		Email:    creds.Email,
		Name:     creds.Name,
		Password: creds.Password,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to register user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}

		return
	}

	userSession := response.SessionId
	jwtToken := response.Jwttoken

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	body, err := utils.ConvertAuthResponseProtoToGo(response, userSession)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert auth response",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if _, err = easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start).Seconds()
	logger.AccessLogger.Info("Completed RegisterUser request",
		zap.String("request_id", requestID),
		zap.Duration("duration", time.Duration(duration)),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received LoginUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	var creds domain.User
	if err = easyjson.UnmarshalFromReader(r.Body, &creds); err != nil {
		logger.AccessLogger.Error("Failed to decode request body",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
		return
	}

	csrfToken, _ := r.Cookie("csrf_token")
	if csrfToken != nil {
		logger.AccessLogger.Error("csrf_token already exists",
			zap.String("request_id", requestID),
			zap.Error(errors.New("csrf_token already exists")),
		)
		err = errors.New("csrf_token already exists")
		statusCode = h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.LoginUser(ctx, &gen.LoginUserRequest{
		Username: creds.Username,
		Password: creds.Password,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to login user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	userSession := response.SessionId
	jwtToken := response.Jwttoken

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    userSession,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	body, err := utils.ConvertAuthResponseProtoToGo(response, userSession)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert auth response",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(body, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received LogoutUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	authHeader := r.Header.Get("X-CSRF-Token")

	_, err = h.client.LogoutUser(ctx, &gen.LogoutRequest{
		AuthHeader: authHeader,
		SessionId:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to logout user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	logoutResponse := domain.ResponseMessage{
		Message: "Successfully logged out",
	}
	if _, err = easyjson.MarshalToWriter(logoutResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()
	logger.AccessLogger.Info("Received PutUser request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	authHeader := r.Header.Get("X-CSRF-Token")

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	metadata := r.FormValue("metadata")

	var creds domain.User
	if err = creds.UnmarshalJSON([]byte(metadata)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.AccessLogger.Warn("Failed to parse metadata",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, errors.New("invalid metadata JSON"), requestID)
		return
	}
	var fileBytes []byte
	var avatar *multipart.FileHeader
	if len(r.MultipartForm.File["avatar"]) > 0 {
		avatar = r.MultipartForm.File["avatar"][0]
		file, err := avatar.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open avatar file",
				zap.String("request_id", requestID),
				zap.Error(err))
			statusCode = h.handleError(w, err, requestID)
			return
		}
		defer file.Close()

		fileBytes, err = io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read avatar file",
				zap.String("request_id", requestID),
				zap.Error(err))
			statusCode = h.handleError(w, err, requestID)
			return
		}
	}

	_, err = h.client.PutUser(ctx, &gen.PutUserRequest{
		Creds: &gen.Metadata{
			Uuid:       creds.UUID,
			Username:   creds.Username,
			Password:   creds.Password,
			Email:      creds.Email,
			Name:       creds.Name,
			Score:      float32(creds.Score),
			Avatar:     creds.Avatar,
			Sex:        creds.Sex,
			GuestCount: int32(creds.GuestCount),
			Birthdate:  timestamppb.New(creds.Birthdate),
			IsHost:     creds.IsHost,
		},
		AuthHeader: authHeader,
		SessionId:  sessionID,
		Avatar:     fileBytes,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to update user data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	PutUserResponse := domain.ResponseMessage{
		Message: "Successfully update user data",
	}
	if _, err = easyjson.MarshalToWriter(PutUserResponse, w); err != nil {
		logger.AccessLogger.Error("Failed to encode update response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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
	userId := mux.Vars(r)["userId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		sanitizedPath := metrics.SanitizeUserIdPath(r.URL.Path)
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, sanitizedPath, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, sanitizedPath, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetUserById request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	user, err := h.client.GetUserById(ctx, &gen.GetUserByIdRequest{
		UserId: userId,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get user by id",
			zap.String("request_id", requestID),
			zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	response, err := utils.ConvertUserResponseProtoToGo(user.User)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert user response",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(response, w); err != nil {
		logger.AccessLogger.Error("Failed to encode getUserById response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetAllUsers request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	users, err := h.client.GetAllUsers(ctx, &gen.Empty{})
	if err != nil {
		logger.AccessLogger.Error("Failed to get all users data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	body, err := utils.ConvertUsersProtoToGo(users)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert users proto",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	response := domain.GetAllUsersResponse{
		Users: body,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = easyjson.MarshalToWriter(response, w); err != nil {
		logger.AccessLogger.Error("Failed to encode getAllUsers response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received GetSessionData request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	sessionData, err := h.client.GetSessionData(ctx, &gen.GetSessionDataRequest{
		SessionId: sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get session data",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	response, err := utils.ConvertSessionDataProtoToGo(sessionData)
	if err != nil {
		logger.AccessLogger.Error("Failed to convert session data",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	if _, err = easyjson.MarshalToWriter(response, w); err != nil {
		logger.AccessLogger.Error("Failed to encode GetSessionData response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
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
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()
	var err error
	statusCode := http.StatusOK
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	logger.AccessLogger.Info("Received RefreshCsrfToken request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}
	newCsrfToken, err := h.client.RefreshCsrfToken(ctx, &gen.RefreshCsrfTokenRequest{
		SessionId: sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to generate CSRF",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    newCsrfToken.CsrfToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	response := domain.CSRFTokenResponse{
		Token: newCsrfToken.CsrfToken,
	}
	if _, err = easyjson.MarshalToWriter(response, w); err != nil {
		logger.AccessLogger.Error("Failed to encode getUserById response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		statusCode = h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed RefreshCsrfToken request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AuthHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := domain.ErrorResponse{
		Error: err.Error(),
	}
	var statusCode int
	switch err.Error() {
	case "username, password, and email are required",
		"username and password are required",
		"invalid credentials",
		"csrf_token already exists",
		"Input contains invalid characters",
		"Input exceeds character limit",
		"Invalid size, type or resolution of image",
		"invalid metadata JSON",
		"missing X-CSRF-Token header",
		"invalid JWT token",
		"invalid type for id in session data",
		"invalid type for avatar in session data":
		statusCode = http.StatusBadRequest

	case "user already exists",
		"email already exists",
		"session already exists",
		"already logged in":
		statusCode = http.StatusConflict

	case "no active session",
		"session not found",
		"user ID not found in session":
		statusCode = http.StatusUnauthorized

	case "user not found",
		"error fetching user by ID",
		"error fetching user by name",
		"error fetching user by email",
		"there is none user in db":
		statusCode = http.StatusNotFound

	case "error creating user",
		"error saving user",
		"error updating user",
		"error fetching all users",
		"failed to generate error response",
		"failed to hash password",
		"failed to upload file",
		"failed to delete file",
		"failed to generate session id",
		"failed to save session",
		"failed to delete session",
		"failed to get user ID",
		"failed to get session data",
		"failed to refresh csrf token",
		"error generating random bytes for session ID",
		"failed to get session id from request cookie",
		"token parse error",
		"token invalid",
		"token expired",
		"bad sign method":
		statusCode = http.StatusInternalServerError

	default:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	if _, jsonErr := easyjson.MarshalToWriter(&errorResponse, w); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
	return statusCode
}
