package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreateUser called", zap.String("request_id", requestID), zap.String("username", creds.Username))

	if err := r.db.Create(creds).Error; err != nil {
		logger.DBLogger.Error("Error creating user", zap.String("request_id", requestID), zap.String("username", creds.Username), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully created user", zap.String("request_id", requestID), zap.String("username", creds.Username))
	return nil
}

func (r *authRepository) SaveUser(ctx context.Context, creds *domain.User) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("SaveUser called", zap.String("request_id", requestID), zap.String("username", creds.Username))
	if err := r.db.Save(creds).Error; err != nil {
		logger.DBLogger.Error("Error saving user", zap.String("request_id", requestID), zap.String("username", creds.Username), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully saved user", zap.String("request_id", requestID), zap.String("username", creds.Username))
	return nil
}

func (r *authRepository) PutUser(ctx context.Context, creds *domain.User, userID string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("PutUser called", zap.String("request_id", requestID), zap.String("userID", userID))

	if err := r.db.Model(&domain.User{}).Where("UUID = ?", userID).Updates(creds).Error; err != nil {
		logger.DBLogger.Error("Error updating user", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
		return err
	}
	//для булевых false
	if err := r.db.Model(&domain.User{}).Where("UUID = ?", userID).Update("isHost", creds.IsHost).Error; err != nil {
		logger.DBLogger.Error("Error updating user", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully updated user", zap.String("request_id", requestID), zap.String("userID", userID))
	return nil
}

func (r *authRepository) GetAllUser(ctx context.Context) ([]domain.User, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAllUser called", zap.String("request_id", requestID))

	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		logger.DBLogger.Error("Error fetching all users", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched all users", zap.String("request_id", requestID), zap.Int("count", len(users)))
	return users, nil
}

func (r *authRepository) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserById called", zap.String("request_id", requestID), zap.String("userID", userID))

	var user domain.User
	if err := r.db.Where("uuid = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found", zap.String("request_id", requestID), zap.String("userID", userID))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by ID", zap.String("request_id", requestID), zap.String("userID", userID), zap.Error(err))
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched user by ID", zap.String("request_id", requestID), zap.String("userID", userID))
	return &user, nil
}

func (r *authRepository) GetUserByName(ctx context.Context, username string) (*domain.User, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserByName called", zap.String("request_id", requestID), zap.String("username", username))

	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.DBLogger.Warn("User not found by username", zap.String("request_id", requestID), zap.String("username", username))
			return nil, errors.New("user not found")
		}
		logger.DBLogger.Error("Error fetching user by name", zap.String("request_id", requestID), zap.String("username", username), zap.Error(err))
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched user by name", zap.String("request_id", requestID), zap.String("username", username))
	return &user, nil
}
