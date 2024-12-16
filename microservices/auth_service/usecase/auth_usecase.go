package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/images"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/validation"
	"context"
	"errors"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

type AuthUseCase interface {
	RegisterUser(ctx context.Context, creds *domain.User) error
	LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error)
	PutUser(ctx context.Context, creds *domain.User, userID string, avatar []byte) error
	GetAllUser(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, userID string) (*domain.User, error)
}

type authUseCase struct {
	authRepository domain.AuthRepository
	minioService   images.MinioServiceInterface
}

func NewAuthUseCase(authRepository domain.AuthRepository, minioService images.MinioServiceInterface) AuthUseCase {
	return &authUseCase{
		authRepository: authRepository,
		minioService:   minioService,
	}
}

func (uc *authUseCase) RegisterUser(ctx context.Context, creds *domain.User) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s]*$`)
	if !validCharPattern.MatchString(creds.Avatar) ||
		!validCharPattern.MatchString(creds.UUID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(creds.Username) > maxLen || len(creds.Email) > maxLen || len(creds.Password) > maxLen || len(creds.Name) > maxLen || len(creds.Avatar) > maxLen || len(creds.UUID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	if creds.Username == "" || creds.Password == "" || creds.Email == "" {
		return errors.New("username, password, and email are required")
	}
	errorResponse := domain.WrongFieldErrorResponse{
		Error:       "Incorrect data forms",
		WrongFields: make([]string, 0),
	}
	var wrongFields []string
	if !validation.ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidateEmail(creds.Email) {
		wrongFields = append(wrongFields, "email")
	}
	if !validation.ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if !validation.ValidateName(creds.Name) {
		wrongFields = append(wrongFields, "name")
	}
	if len(wrongFields) > 0 {
		errorResponse.WrongFields = wrongFields
		errorResponseJSON, err := easyjson.Marshal(errorResponse)
		if err != nil {
			return errors.New("failed to generate error response")
		}
		return errors.New(string(errorResponseJSON))
	}

	// Хэширование пароля
	hashedPassword, ok := middleware.HashPassword(creds.Password)
	if ok != nil {
		return errors.New("failed to hash password")
	}
	creds.Password = hashedPassword

	existingUser, _ := uc.authRepository.GetUserByName(ctx, creds.Username)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	existingUser, _ = uc.authRepository.GetUserByEmail(ctx, creds.Email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	err := uc.authRepository.CreateUser(ctx, creds)
	if err != nil {
		return err
	}

	return uc.authRepository.SaveUser(ctx, creds)
}

func (uc *authUseCase) LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error) {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s]*$`)
	if !validCharPattern.MatchString(creds.Email) ||
		!validCharPattern.MatchString(creds.Name) ||
		!validCharPattern.MatchString(creds.Avatar) ||
		!validCharPattern.MatchString(creds.UUID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return nil, errors.New("input contains invalid characters")
	}

	if len(creds.Username) > maxLen || len(creds.Email) > maxLen || len(creds.Password) > maxLen || len(creds.Name) > maxLen || len(creds.Avatar) > maxLen || len(creds.UUID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return nil, errors.New("input exceeds character limit")
	}

	if creds.Username == "" || creds.Password == "" {
		return nil, errors.New("username and password are required")
	}
	errorResponse := domain.WrongFieldErrorResponse{
		Error:       "Incorrect data forms",
		WrongFields: make([]string, 0),
	}
	var wrongFields []string
	if !validation.ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if len(wrongFields) > 0 {
		errorResponse.WrongFields = wrongFields
		errorResponseJSON, err := easyjson.Marshal(errorResponse)
		if err != nil {
			return nil, errors.New("failed to generate error response")
		}
		return nil, errors.New(string(errorResponseJSON))
	}

	requestedUser, err := uc.authRepository.GetUserByName(ctx, creds.Username)
	if err != nil || requestedUser == nil {
		return nil, errors.New("user not found")
	}

	if !middleware.CheckPassword(requestedUser.Password, creds.Password) {
		return nil, errors.New("invalid credentials")
	}

	return requestedUser, nil
}

func (uc *authUseCase) PutUser(ctx context.Context, creds *domain.User, userID string, avatar []byte) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s\-_]*$`)
	if !validCharPattern.MatchString(creds.Avatar) ||
		!validCharPattern.MatchString(creds.UUID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(creds.Username) > maxLen || len(creds.Email) > maxLen || len(creds.Password) > maxLen || len(creds.Name) > maxLen || len(creds.Avatar) > maxLen || len(creds.UUID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	if avatar != nil {
		if err := validation.ValidateImage(avatar, 5<<20, []string{"image/jpeg", "image/png", "image/jpg"}, 2000, 2000); err != nil {
			logger.AccessLogger.Warn("Invalid size, type or resolution of image", zap.String("request_id", requestID), zap.Error(err))
			return err
		}
	}

	var wrongFields []string
	errorResponse := domain.WrongFieldErrorResponse{
		Error:       "Incorrect data forms",
		WrongFields: make([]string, 0),
	}
	if !validation.ValidateLogin(creds.Username) && len(creds.Username) > 0 {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidateEmail(creds.Email) && len(creds.Email) > 0 {
		wrongFields = append(wrongFields, "email")
	}
	if !validation.ValidatePassword(creds.Password) && len(creds.Password) > 0 {
		wrongFields = append(wrongFields, "password")
	}
	if !validation.ValidateName(creds.Name) && len(creds.Name) > 0 {
		wrongFields = append(wrongFields, "name")
	}
	if len(wrongFields) > 0 {
		errorResponse.WrongFields = wrongFields
		errorResponseJSON, err := easyjson.Marshal(errorResponse)
		if err != nil {
			return errors.New("failed to generate error response")
		}
		return errors.New(string(errorResponseJSON))
	}

	if avatar != nil {
		contentType := http.DetectContentType(avatar[:512])

		uploadedPath, err := uc.minioService.UploadFile(avatar, contentType, "user/"+userID)
		if err != nil {
			logger.AccessLogger.Warn("Failed to upload file", zap.String("request_id", requestID), zap.Error(err))
			return errors.New("failed to upload file")
		}

		creds.Avatar = "/images/" + uploadedPath
	}

	err := uc.authRepository.PutUser(ctx, creds, userID)
	if err != nil {
		ok := uc.minioService.DeleteFile(creds.Avatar)
		if ok != nil {
			return errors.New("failed to delete file")
		}
		return err
	}
	return nil
}

func (uc *authUseCase) GetAllUser(ctx context.Context) ([]domain.User, error) {
	users, err := uc.authRepository.GetAllUser(ctx)
	if err != nil {
		return nil, errors.New("there is none user in db")
	}
	return users, nil
}

func (uc *authUseCase) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(userID) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return nil, errors.New("input contains invalid characters")
	}

	if len(userID) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return nil, errors.New("input exceeds character limit")
	}

	user, err := uc.authRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
