package session

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"net/http"
)

type ServiceSession struct {
	store *sessions.CookieStore
}

func NewSessionService(store *sessions.CookieStore) *ServiceSession {
	return &ServiceSession{store: store}
}

func (s *ServiceSession) LogoutSession(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("LogoutSession called", zap.String("request_id", requestID))

	session, err := s.store.Get(r, "session_id")
	if err != nil {
		logger.AccessLogger.Error("Error retrieving session", zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	if session.IsNew {
		logger.AccessLogger.Warn("Attempted to logout with no active session", zap.String("request_id", requestID))
		return errors.New("no active session")
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		logger.AccessLogger.Error("Error saving session during logout", zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	logger.AccessLogger.Info("Successfully logged out session", zap.String("request_id", requestID))
	return nil
}

func (s *ServiceSession) GetUserID(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("GetUserID called", zap.String("request_id", requestID))

	session, err := s.store.Get(r, "session_id")
	if err != nil {
		logger.AccessLogger.Error("Error retrieving session", zap.String("request_id", requestID), zap.Error(err))
		return "", err
	}

	if session.IsNew {
		logger.AccessLogger.Warn("No active session found when getting user ID", zap.String("request_id", requestID))
		return "", errors.New("no active session")
	}

	userID, ok := session.Values["id"].(string)
	if !ok {
		logger.AccessLogger.Error("User ID not found or invalid type in session", zap.String("request_id", requestID))
		return "", errors.New("user ID not found in session")
	}

	logger.AccessLogger.Info("Successfully retrieved user ID from session", zap.String("request_id", requestID), zap.String("userID", userID))
	return userID, nil
}

func (s *ServiceSession) GetSessionData(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("GetSessionData called", zap.String("request_id", requestID))

	session, err := s.store.Get(r, "session_id")
	if err != nil {
		logger.AccessLogger.Error("Error retrieving session", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	if session.IsNew {
		logger.AccessLogger.Warn("No active session found when getting session data", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	userID, okID := session.Values["id"].(string)
	avatar, okAvatar := session.Values["avatar"].(string)

	if !okID {
		logger.AccessLogger.Error("User ID not found or invalid type in session", zap.String("request_id", requestID))
		return nil, errors.New("user ID not found in session")
	}

	sessionData := map[string]interface{}{
		"id":     userID,
		"avatar": "",
	}

	if okAvatar {
		sessionData["avatar"] = avatar
	}

	logger.AccessLogger.Info("Successfully retrieved session data", zap.String("request_id", requestID), zap.Any("sessionData", sessionData))
	return &sessionData, nil
}

func (s *ServiceSession) CreateSession(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("CreateSession called", zap.String("request_id", requestID), zap.String("userID", user.UUID))

	session, err := s.store.Get(r, "session_id")
	if err != nil {
		logger.AccessLogger.Error("Error retrieving session for creation", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	if !session.IsNew {
		logger.AccessLogger.Warn("Attempted to create a session when one already exists", zap.String("request_id", requestID), zap.String("userID", user.UUID))
		return nil, errors.New("session already exists")
	}

	session.Values["id"] = user.UUID
	session.Values["username"] = user.Username
	session.Values["email"] = user.Email

	if user.Name != "" {
		session.Values["name"] = user.Name
	}
	if user.Avatar != "" {
		session.Values["avatar"] = user.Avatar
	}

	sessionID, err := GenerateSessionID(ctx)
	if err != nil {
		logger.AccessLogger.Error("Failed to generate session ID", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("failed to generate session id")
	}

	session.Values["session_id"] = sessionID

	if err := session.Save(r, w); err != nil {
		logger.AccessLogger.Error("Failed to save session", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("failed to save session")
	}

	logger.AccessLogger.Info("Successfully created session", zap.String("request_id", requestID), zap.String("session_id", sessionID), zap.String("userID", user.UUID))
	return session, nil
}

func GenerateSessionID(ctx context.Context) (string, error) {
	// Включаем логирование при генерации session ID
	requestID := middleware.GetRequestID(ctx)
	logger.AccessLogger.Info("GenerateSessionID called", zap.String("request_id", requestID))

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		logger.AccessLogger.Error("Error generating random bytes for session ID", zap.String("request_id", requestID), zap.Error(err))
		return "", err
	}

	sessionID := base64.StdEncoding.EncodeToString(b)
	logger.AccessLogger.Info("Successfully generated session ID", zap.String("request_id", requestID), zap.String("session_id", sessionID))
	return sessionID, nil
}
