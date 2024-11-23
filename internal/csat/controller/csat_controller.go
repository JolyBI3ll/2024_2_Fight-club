package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/microservices/csat_service/controller/gen"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type CsatHandler struct {
	client         gen.CsatClient
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewCsatHandler(client gen.CsatClient, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *CsatHandler {
	return &CsatHandler{
		client:         client,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *CsatHandler) GetServey(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	sanitizer := bluemonday.UGCPolicy()
	ctx, cancel := middleware.WithTimeout(r.Context())
	surveyId, err := strconv.Atoi(sanitizer.Sanitize(mux.Vars(r)["surveyId"]))
	if err != nil {
		logger.AccessLogger.Error("Failed to atoi",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}
	defer cancel()

	logger.AccessLogger.Info("Received GetSurvey request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	_, err = h.sessionService.GetUserID(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Error("Failed to get user ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.GetSurvey(ctx, &gen.GetSurveyRequest{
		SurveyId: int32(surveyId),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get survey",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode update response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetSurvey request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *CsatHandler) PostAnswer(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received PostAnswer request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	authHeader := r.Header.Get("X-CSRF-Token")

	var answers []*gen.Answer
	if err := json.NewDecoder(r.Body).Decode(&answers); err != nil {
		logger.AccessLogger.Error("Failed to decode request body",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, errors.New("failed to decode"), requestID)
		return
	}

	response, err := h.client.PostAnswers(ctx, &gen.PostAnswersRequest{
		SessionId:  sessionID,
		AuthHeader: authHeader,
		Answer:     answers,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to post answers",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response.Response); err != nil {
		logger.AccessLogger.Error("Failed to encode update response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetSurvey request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *CsatHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	logger.AccessLogger.Info("Received GetStatistics request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}
	_, err = h.sessionService.GetUserID(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Error("Failed to get user ID",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	grpcResponse, err := h.client.GetStatistics(ctx, &gen.Empty{})
	if err != nil {
		logger.AccessLogger.Warn("Failed to get statistics",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	var response []domain.GetStatictics
	for _, grpcStat := range grpcResponse.Statistics {
		answerNumbers := make(map[int]int)
		for _, item := range grpcStat.Map {
			answerNumbers[int(item.Key)] = int(item.Values)
		}
		response = append(response, domain.GetStatictics{
			Avg:           grpcStat.Average,
			Title:         grpcStat.Title,
			AnswerNumbers: answerNumbers,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode statistics response",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetStatistics request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *CsatHandler) handleError(w http.ResponseWriter, err error, requestID string) {
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
		"session already exists",
		"email already exists":
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
