package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/auth/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{}
	mockMinioService := &mocks.MockMinioService{}

	uc := NewAuthUseCase(mockAuthRepo, mockMinioService)
	ctx := context.TODO()

	creds := &domain.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Name:     "Test User",
	}

	// Тест-кейс 1: Успешная регистрация без аватара
	mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
		return nil, nil
	}
	mockAuthRepo.CreateUserFunc = func(ctx context.Context, user *domain.User) error {
		return nil
	}
	mockAuthRepo.SaveUserFunc = func(ctx context.Context, user *domain.User) error {
		return nil
	}

	err := uc.RegisterUser(ctx, creds)
	assert.NoError(t, err)

	// Тест-кейс 2: Ошибка валидации - неправильный логин
	invalidCreds := &domain.User{
		Username: "t", // некорректный логин
		Password: "password",
		Email:    "test@example.com",
		Name:     "Test User",
	}

	err = uc.RegisterUser(ctx, invalidCreds)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username")

	// Тест-кейс 3: Ошибка - пользователь уже существует
	mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
		return creds, nil
	}

	err = uc.RegisterUser(ctx, creds)
	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
}

func TestLoginUser(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{}
	uc := NewAuthUseCase(mockAuthRepo, nil)
	ctx := context.TODO()

	creds := &domain.User{
		Username: "testuser",
		Password: "password",
	}

	// Тест-кейс 1: Успешный вход
	mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
		return creds, nil
	}

	user, err := uc.LoginUser(ctx, creds)
	assert.NoError(t, err)
	assert.Equal(t, creds, user)

	// Тест-кейс 2: Неправильные учетные данные
	invalidCreds := &domain.User{
		Username: "testuser",
		Password: "wrongpassword",
	}
	user, err = uc.LoginUser(ctx, invalidCreds)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestPutUser(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{}
	mockMinioService := &mocks.MockMinioService{}
	uc := NewAuthUseCase(mockAuthRepo, mockMinioService)
	ctx := context.TODO()

	creds := &domain.User{
		Username: "testuser",
		Password: "password",
	}
	userID := "test-uuid"

	// Тест-кейс 1: Успешное обновление пользователя без аватара
	mockAuthRepo.PutUserFunc = func(ctx context.Context, user *domain.User, userID string) error {
		return nil
	}

	err := uc.PutUser(ctx, creds, userID, nil)
	assert.NoError(t, err)

	// Тест-кейс 2: Ошибка при загрузке аватара
	mockMinioService.UploadFileFunc = func(file *multipart.FileHeader, path string) (string, error) {
		return "", errors.New("upload error")
	}

	avatar := &multipart.FileHeader{Filename: "avatar.jpg"}
	err = uc.PutUser(ctx, creds, userID, avatar)
	assert.Error(t, err)
	assert.Equal(t, "upload error", err.Error())
}

func TestGetAllUser(t *testing.T) {
	mockAuthRepo := &mocks.MockAuthRepository{}
	uc := NewAuthUseCase(mockAuthRepo, nil)
	ctx := context.TODO()

	// Тест-кейс 1: Успешное получение всех пользователей
	mockAuthRepo.GetAllUserFunc = func(ctx context.Context) ([]domain.User, error) {
		return []domain.User{{Username: "testuser"}}, nil
	}

	users, err := uc.GetAllUser(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, users)

	// Тест-кейс 2: Ошибка - нет пользователей
	mockAuthRepo.GetAllUserFunc = func(ctx context.Context) ([]domain.User, error) {
		return nil, errors.New("no users")
	}

	users, err = uc.GetAllUser(ctx)
	assert.Error(t, err)
	assert.Nil(t, users)
}
