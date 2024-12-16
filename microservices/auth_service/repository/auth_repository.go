package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, creds *domain.User) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreateUser called", zap.String("request_id", requestID), zap.String("username", creds.Username))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("CreateUser", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("CreateUser", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("CreateUser").Observe(duration)
	}()
	if err := r.db.Create(creds).Error; err != nil {
		logger.DBLogger.Error("Error creating user", zap.String("request_id", requestID), zap.String("username", creds.Username), zap.Error(err))
		return errors.New("error creating user")
	}

	logger.DBLogger.Info("Successfully created user", zap.String("request_id", requestID), zap.String("username", creds.Username))
	return nil
}

func (r *authRepository) SaveUser(ctx context.Context, creds *domain.User) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("SaveUser called", zap.String("request_id", requestID), zap.String("username", creds.Username))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("SaveUser", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("SaveUser", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("SaveUser").Observe(duration)
	}()
	if err := r.db.Save(creds).Error; err != nil {
		logger.DBLogger.Error("Error saving user", zap.String("request_id", requestID), zap.String("username", creds.Username), zap.Error(err))
		return errors.New("error saving user")
	}

	logger.DBLogger.Info("Successfully saved user", zap.String("request_id", requestID), zap.String("username", creds.Username))
	return nil
}

func (r *authRepository) PutUser(ctx context.Context, creds *domain.User, userID string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("PutUser called", zap.String("request_id", requestID), zap.String("userID", userID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("PutUser", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("PutUser", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("PutUser").Observe(duration)
	}()

	if err := r.db.Model(&domain.User{}).Where("UUID = ?", userID).Updates(creds).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			logger.DBLogger.Warn("Unique constraint violation", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
			return errors.New("username or email already exists")
		}
		logger.DBLogger.Error("Error updating user", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
		return errors.New("error updating user")
	}

	if err := r.db.Model(&domain.User{}).Where("UUID = ?", userID).Update("isHost", creds.IsHost).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			logger.DBLogger.Warn("Unique constraint violation on isHost", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
			return errors.New("username or email already exists")
		}
		logger.DBLogger.Error("Error updating user", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
		return errors.New("error updating user")
	}

	logger.DBLogger.Info("Successfully updated user", zap.String("request_id", requestID), zap.String("userID", userID))
	return nil
}

func (r *authRepository) GetAllUser(ctx context.Context) ([]domain.User, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAllUser called", zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetAllUser", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetAllUser", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetAllUser").Observe(duration)
	}()
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		logger.DBLogger.Error("Error fetching all users", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching all users")
	}

	logger.DBLogger.Info("Successfully fetched all users", zap.String("request_id", requestID), zap.Int("count", len(users)))
	return users, nil
}

func (r *authRepository) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserById called", zap.String("request_id", requestID), zap.String("userID", userID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetUserById", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetUserById", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetUserById").Observe(duration)
	}()
	var user domain.User
	if err := r.db.Where("uuid = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found", zap.String("request_id", requestID), zap.String("userID", userID))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by ID", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
		return nil, errors.New("error fetching user by ID")
	}

	logger.DBLogger.Info("Successfully fetched user by ID", zap.String("request_id", requestID), zap.String("userID", userID))
	return &user, nil
}

func (r *authRepository) GetUserByName(ctx context.Context, username string) (*domain.User, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserByName called", zap.String("request_id", requestID), zap.String("username", username))
	start := time.Now()
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetUserByName", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetUserByName", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetUserByName").Observe(duration)
	}()
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found by username", zap.String("request_id", requestID), zap.String("username", username))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by name", zap.String("request_id", requestID), zap.String("username", username), zap.Error(err))
		return nil, errors.New("error fetching user by name")
	}

	logger.DBLogger.Info("Successfully fetched user by name", zap.String("request_id", requestID), zap.String("username", username))
	return &user, nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserByEmail called", zap.String("request_id", requestID), zap.String("email", email))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetUserByEmail", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetUserByEmail", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetUserByEmail").Observe(duration)
	}()
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found by email", zap.String("request_id", requestID), zap.String("email", email))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by email", zap.String("request_id", requestID), zap.String("email", email), zap.Error(err))
		return nil, errors.New("error fetching user by email")
	}

	logger.DBLogger.Info("Successfully fetched user by email", zap.String("request_id", requestID), zap.String("email", email))
	return &user, nil
}
