package session

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type InterfaceSession interface {
	GetUserID(ctx context.Context, sessionID string) (string, error)
	LogoutSession(ctx context.Context, sessionID string) error
	CreateSession(ctx context.Context, user *domain.User) (string, error)
	GetSessionData(ctx context.Context, sessionID string) (*map[string]interface{}, error)
}

type ServiceSession struct {
	store RedisInterface
}

func NewSessionService(store RedisInterface) InterfaceSession {
	return &ServiceSession{
		store: store,
	}
}

// CreateSession Создание сессии
func (s *ServiceSession) CreateSession(ctx context.Context, user *domain.User) (string, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("CreateSession called", zap.String("request_id", requestID), zap.String("userID", user.UUID))

	// Генерация уникального session_id
	sessionID, err := GenerateSessionID(ctx)
	if err != nil {
		logger.AccessLogger.Error("Failed to generate session ID", zap.String("request_id", requestID), zap.Error(err))
		return "", errors.New("failed to generate session id")
	}

	// Данные для сессии
	sessionData := map[string]interface{}{
		"id":     user.UUID,
		"avatar": user.Avatar,
	}

	// Сохранение сессии в Redis
	if err := s.store.Set(ctx, sessionID, sessionData, 24*time.Hour); err != nil {
		logger.AccessLogger.Error("Failed to save session", zap.String("request_id", requestID), zap.Error(err))
		return "", errors.New("failed to save session")
	}

	logger.AccessLogger.Info("Successfully created session", zap.String("request_id", requestID), zap.String("session_id", sessionID), zap.String("userID", user.UUID))
	return sessionID, nil
}

// GetUserID Получение данных пользователя по session_id
func (s *ServiceSession) GetUserID(ctx context.Context, sessionID string) (string, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("GetUserID called", zap.String("request_id", requestID))

	data, err := s.store.Get(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session", zap.String("request_id", requestID), zap.Error(err))
		return "", errors.New("session not found")
	}

	userID, ok := data["id"].(string)
	if !ok {
		logger.AccessLogger.Error("User ID not found or invalid type in session", zap.String("request_id", requestID))
		return "", errors.New("user ID not found in session")
	}

	logger.AccessLogger.Info("Successfully retrieved user ID from session", zap.String("request_id", requestID), zap.String("userID", userID))
	return userID, nil
}

// GetSessionData Получение данных сессии
func (s *ServiceSession) GetSessionData(ctx context.Context, sessionID string) (*map[string]interface{}, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("GetSessionData called", zap.String("request_id", requestID))

	data, err := s.store.Get(ctx, sessionID)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("session not found")
	}

	logger.AccessLogger.Info("Successfully retrieved session data", zap.String("request_id", requestID), zap.Any("session_data", data))
	return &data, nil
}

// LogoutSession Удаление сессии
func (s *ServiceSession) LogoutSession(ctx context.Context, sessionID string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("LogoutSession called", zap.String("request_id", requestID))

	if err := s.store.Delete(ctx, sessionID); err != nil {
		logger.AccessLogger.Error("Failed to delete session", zap.String("request_id", requestID), zap.Error(err))
		return errors.New("failed to delete session")
	}

	logger.AccessLogger.Info("Successfully logged out session", zap.String("request_id", requestID))
	return nil
}

// GenerateSessionID Генерация уникального session_id
func GenerateSessionID(ctx context.Context) (string, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("GenerateSessionID called", zap.String("request_id", requestID))

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logger.AccessLogger.Error("Error generating random bytes for session ID", zap.String("request_id", requestID), zap.Error(err))
		return "", errors.New("error generating random bytes for session ID")
	}

	sessionID := base64.StdEncoding.EncodeToString(b)
	logger.AccessLogger.Info("Successfully generated session ID", zap.String("request_id", requestID), zap.String("session_id", sessionID))
	return sessionID, nil
}

func GetSessionId(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", errors.New("failed to get session id from request cookie")
	}

	sessionID := cookie.Value
	return sessionID, nil
}
