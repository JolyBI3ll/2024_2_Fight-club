package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/validation"
	"context"
	"encoding/json"
	"errors"
)

type AuthUseCase interface {
	RegisterUser(ctx context.Context, creds *domain.User) error
	LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error)
	PutUser(ctx context.Context, creds *domain.User, userID string) error
	GetAllUser(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, userID string) (*domain.User, error)
}

type authUseCase struct {
	authRepository domain.AuthRepository
}

func NewAuthUseCase(authRepository domain.AuthRepository) AuthUseCase {
	return &authUseCase{
		authRepository: authRepository,
	}
}

func (uc *authUseCase) RegisterUser(ctx context.Context, creds *domain.User) error {
	if creds.Username == "" || creds.Password == "" || creds.Email == "" {
		return errors.New("username, password, and email are required")
	}
	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
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
		errorResponse["wrongFields"] = wrongFields
		errorResponseJSON, err := json.Marshal(errorResponse)
		if err != nil {
			return errors.New("failed to generate error response")
		}
		return errors.New(string(errorResponseJSON))
	}
	existingUser, _ := uc.authRepository.GetUserByName(ctx, creds.Username)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	return uc.authRepository.CreateUser(ctx, creds)
}

func (uc *authUseCase) LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error) {
	if creds.Username == "" || creds.Password == "" {
		return nil, errors.New("username and password are required")
	}
	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
	}
	var wrongFields []string
	if !validation.ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if len(wrongFields) > 0 {
		errorResponse["wrongFields"] = wrongFields
		errorResponseJSON, err := json.Marshal(errorResponse)
		if err != nil {
			return nil, errors.New("failed to generate error response")
		}
		return nil, errors.New(string(errorResponseJSON))
	}
	requestedUser, err := uc.authRepository.GetUserByName(ctx, creds.Username)
	if err != nil || requestedUser == nil {
		return nil, errors.New("user not found")
	}
	if requestedUser.Password != creds.Password {
		return nil, errors.New("invalid credentials")
	}
	return requestedUser, nil
}

func (uc *authUseCase) PutUser(ctx context.Context, creds *domain.User, userID string) error {
	err := uc.authRepository.PutUser(ctx, creds, userID)
	if err != nil {
		return errors.New("user not found")
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
	user, err := uc.authRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
