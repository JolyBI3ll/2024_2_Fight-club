package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/microservices/auth_service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"testing"
)

// GenerateImage создает изображение с указанным форматом и размером
func GenerateImage(format string, width, height int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Заливаем белым цветом
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, white)
		}
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported format")
	}

	return buf.Bytes(), nil
}

func TestRegisterUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockAuthRepo := &mocks.MockAuthRepository{}
	mockMinioService := &mocks.MockMinioService{}

	uc := NewAuthUseCase(mockAuthRepo, mockMinioService)
	ctx := context.TODO()

	// Тест-кейс 1: Успешная регистрация
	t.Run("Successful Registration", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.MockGetUserByEmail = func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.CreateUserFunc = func(ctx context.Context, user *domain.User) error {
			return nil
		}
		mockAuthRepo.SaveUserFunc = func(ctx context.Context, user *domain.User) error {
			return nil
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "avatar.png",
			UUID:     "uuid123",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.NoError(t, err)
	})

	// Тест-кейс 2: Неправильные символы в Avatar или UUID
	t.Run("Invalid Characters", func(t *testing.T) {
		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "invalid<>", // Некорректный символ
			UUID:     "valid-uuid",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input contains invalid characters")
	})

	// Тест-кейс 3: Превышение длины полей
	t.Run("Input Exceeds Length", func(t *testing.T) {
		longString := make([]byte, 256)
		for i := range longString {
			longString[i] = 'a'
		}

		creds := &domain.User{
			Username: string(longString),
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "avatar.png",
			UUID:     "uuid123",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input exceeds character limit")
	})

	// Тест-кейс 4: Отсутствие обязательных полей
	t.Run("Missing Required Fields", func(t *testing.T) {
		creds := &domain.User{
			Username: "",
			Password: "",
			Email:    "",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username, password, and email are required")
	})

	// Тест-кейс 5: Ошибки валидации полей
	t.Run("Validation Errors", func(t *testing.T) {
		creds := &domain.User{
			Username: "u",             // Слишком короткий логин
			Password: "123",           // Слишком простой пароль
			Email:    "invalid_email", // Некорректный email
			Name:     "Invalid<>Name", // Некорректное имя
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username")
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "password")
		assert.Contains(t, err.Error(), "name")
	})

	// Тест-кейс 6: Пользователь уже существует
	t.Run("User Already Exists", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return &domain.User{}, nil
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
	})

	// Тест-кейс 7: Email уже зарегистрирован
	t.Run("Email Already Exists", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.MockGetUserByEmail = func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{}, nil
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
	})

	// Тест-кейс 8: Ошибка сохранения пользователя
	t.Run("User Save Error", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.MockGetUserByEmail = func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.CreateUserFunc = func(ctx context.Context, user *domain.User) error {
			return nil
		}
		mockAuthRepo.SaveUserFunc = func(ctx context.Context, user *domain.User) error {
			return errors.New("save error")
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Equal(t, "save error", err.Error())
	})

	//Тест-кейс 9 ошибка создания пользователя
	t.Run("Create user error", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.MockGetUserByEmail = func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		}
		mockAuthRepo.SaveUserFunc = func(ctx context.Context, user *domain.User) error {
			return nil
		}
		mockAuthRepo.CreateUserFunc = func(ctx context.Context, user *domain.User) error {
			return errors.New("create user error")
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Equal(t, "create user error", err.Error())
	})
}

func TestLoginUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	mockAuthRepo := &mocks.MockAuthRepository{}

	uc := NewAuthUseCase(mockAuthRepo, nil)
	ctx := context.TODO()

	// Тест-кейс 1: Успешный вход
	t.Run("Successful Login", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			password, _ := middleware.HashPassword("password")
			return &domain.User{
				Username: "testuser",
				Password: password,
			}, nil
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
	})

	// Тест-кейс 2: Неправильные символы в Avatar или UUID
	t.Run("Invalid Characters", func(t *testing.T) {
		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "invalid<>", // Некорректный символ
			UUID:     "valid-uuid",
		}

		err := uc.RegisterUser(ctx, creds)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input contains invalid characters")
	})

	// Тест-кейс 3: Некорректные символы в полях
	t.Run("Invalid Characters", func(t *testing.T) {
		creds := &domain.User{
			Username: "invalid<>username#$%",
			Password: "invalid<>password$%^",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "{\"error\":\"Incorrect data forms\",\"wrongFields\":[\"username\",\"password\"]}")
	})

	// Тест-кейс 4: Превышение длины полей
	t.Run("Input Exceeds Length", func(t *testing.T) {
		longString := make([]byte, 512)
		for i := range longString {
			longString[i] = 'a'
		}

		creds := &domain.User{
			Username: string(longString),
			Password: "password",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "input exceeds character limit")
	})

	// Тест-кейс 5: Отсутствие обязательных полей
	t.Run("Missing Required Fields", func(t *testing.T) {
		creds := &domain.User{
			Username: "",
			Password: "",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "username and password are required")
	})

	// Тест-кейс 6: Ошибки валидации полей
	t.Run("Validation Errors", func(t *testing.T) {
		creds := &domain.User{
			Username: "u",     // Слишком короткий логин
			Password: "12345", // Слишком простой пароль
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "username")
		assert.Contains(t, err.Error(), "password")
	})

	// Тест-кейс 7: Пользователь не найден
	t.Run("User Not Found", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return nil, nil
		}

		creds := &domain.User{
			Username: "nonexistent",
			Password: "password",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})

	// Тест-кейс 8: Неверный пароль
	t.Run("Invalid Credentials", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			password, _ := middleware.HashPassword("password")
			return &domain.User{
				Username: "testuser",
				Password: password,
			}, nil
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "wrongpassword",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	// Тест-кейс 9: Ошибка получения пользователя из базы данных
	t.Run("Database Error", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			return nil, errors.New("database error")
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})

	// Тест-кейс 2: Неправильные символы в Avatar или UUID
	t.Run("Invalid Characters", func(t *testing.T) {
		mockAuthRepo.GetUserByNameFunc = func(ctx context.Context, username string) (*domain.User, error) {
			password, _ := middleware.HashPassword("password")
			return &domain.User{
				Username: "testuser",
				Password: password,
			}, nil
		}

		creds := &domain.User{
			Username: "testuser",
			Password: "password",
			Email:    "test@example.com",
			Name:     "Test User",
			Avatar:   "invalid<>", // Некорректный символ
			UUID:     "valid-uuid",
		}

		user, err := uc.LoginUser(ctx, creds)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "input contains invalid characters")
	})
}

func TestPutUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockAuthRepo := &mocks.MockAuthRepository{}
	mockMinioService := &mocks.MockMinioService{}

	uc := NewAuthUseCase(mockAuthRepo, mockMinioService)
	ctx := context.TODO()

	validAvatar, _ := GenerateImage("jpeg", 2000, 2000)
	invalidAvatar, _ := GenerateImage("jpeg", 5000, 5000)

	creds := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "ValidPassword123",
		Name:     "Test User",
		Avatar:   "",
		UUID:     "valid-uuid",
	}

	userID := "test-uuid"

	// Тест-кейс 1: Успешное обновление пользователя без аватара
	t.Run("Successful Update Without Avatar", func(t *testing.T) {
		mockAuthRepo.PutUserFunc = func(ctx context.Context, user *domain.User, userID string) error {
			return nil
		}

		err := uc.PutUser(ctx, creds, userID, nil)
		assert.NoError(t, err)
	})

	// Тест-кейс 2: Успешное обновление пользователя с аватаром
	t.Run("Successful Update With Avatar", func(t *testing.T) {
		mockAuthRepo.PutUserFunc = func(ctx context.Context, user *domain.User, userID string) error {
			return nil
		}
		mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
			return "path/to/image.jpg", nil
		}

		err := uc.PutUser(ctx, creds, userID, validAvatar)
		assert.NoError(t, err)
		assert.Equal(t, "/images/path/to/image.jpg", creds.Avatar)
	})

	// Тест-кейс 3: Неверные символы в данных пользователя
	t.Run("Invalid Characters in User Data", func(t *testing.T) {
		invalidCreds := &domain.User{
			UUID:   "invalid<>uuid",
			Avatar: "invalid<>avatar",
		}

		err := uc.PutUser(ctx, invalidCreds, userID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input contains invalid characters")
	})

	// Тест-кейс 4: Превышение длины данных
	t.Run("Exceeds Character Limit", func(t *testing.T) {
		longString := make([]byte, 256)
		for i := range longString {
			longString[i] = 'a'
		}

		invalidCreds := &domain.User{
			Username: string(longString),
		}

		err := uc.PutUser(ctx, invalidCreds, userID, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input exceeds character limit")
	})

	// Тест-кейс 5: Неверный формат аватара
	t.Run("Invalid Avatar Format", func(t *testing.T) {
		mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
			return "", nil
		}

		creds = &domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "ValidPassword123",
			Name:     "Test User",
			Avatar:   "",
			UUID:     "valid-uuid",
		}

		err := uc.PutUser(ctx, creds, userID, invalidAvatar)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid size, type or resolution of image")
	})

	// Тест-кейс 6: Ошибка загрузки аватара
	t.Run("Avatar Upload Failure", func(t *testing.T) {
		mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
			return "", errors.New("upload error")
		}

		err := uc.PutUser(ctx, creds, userID, validAvatar)
		assert.Error(t, err)
		assert.Equal(t, "failed to upload file", err.Error())
	})

	// Тест-кейс 7: Ошибка удаления файла при откате
	t.Run("Avatar Deletion Failure on Rollback", func(t *testing.T) {
		mockAuthRepo.PutUserFunc = func(ctx context.Context, user *domain.User, userID string) error {
			return errors.New("db error")
		}
		mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
			return "path/to/image.jpg", nil
		}
		mockMinioService.DeleteFileFunc = func(filePath string) error {
			return errors.New("deletion error")
		}

		err := uc.PutUser(ctx, creds, userID, validAvatar)
		assert.Error(t, err)
		assert.Equal(t, "failed to delete file", err.Error())
	})

	// Тест-кейс 8: Ошибка обновления пользователя в базе данных
	t.Run("Database Update Failure", func(t *testing.T) {
		mockAuthRepo.PutUserFunc = func(ctx context.Context, user *domain.User, userID string) error {
			return errors.New("db error")
		}
		mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
			return "path/to/image.jpg", nil
		}
		mockMinioService.DeleteFileFunc = func(filePath string) error {
			return nil
		}

		creds = &domain.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "ValidPassword123",
			Name:     "Test User",
			Avatar:   "",
			UUID:     "valid-uuid",
		}

		err := uc.PutUser(ctx, creds, userID, validAvatar)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	// Тест-кейс 9: Ошибки валидации полей
	t.Run("Validation Errors", func(t *testing.T) {
		creds = &domain.User{
			Username: "u",             // Слишком короткий логин
			Password: "123",           // Слишком простой пароль
			Email:    "invalid_email", // Некорректный email
			Name:     "Invalid<>Name", // Некорректное имя
		}

		err := uc.PutUser(ctx, creds, userID, validAvatar)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username")
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "password")
		assert.Contains(t, err.Error(), "name")
	})
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

func TestGetUserById(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	mockAuthRepo := &mocks.MockAuthRepository{}
	uc := NewAuthUseCase(mockAuthRepo, nil)
	ctx := context.TODO()

	// Успешный тест-кейс
	t.Run("Success", func(t *testing.T) {
		userID := "validUserID"
		expectedUser := &domain.User{
			UUID:     userID,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockAuthRepo.GetUserByIdFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return expectedUser, nil
		}

		user, err := uc.GetUserById(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	// Недопустимые символы в userID
	t.Run("Invalid Characters in UserID", func(t *testing.T) {
		userID := "invalid@ID!"
		_, err := uc.GetUserById(ctx, userID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "input contains invalid characters")
	})

	// Превышение максимальной длины userID
	t.Run("UserID Exceeds Max Length", func(t *testing.T) {
		longString := make([]byte, 256)
		for i := range longString {
			longString[i] = 'a'
		}
		_, err := uc.GetUserById(ctx, string(longString))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "input exceeds character limit")
	})

	// Пользователь не найден
	t.Run("User Not Found", func(t *testing.T) {
		userID := "nonExistentUserID"
		mockAuthRepo.GetUserByIdFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, errors.New("user not found")
		}

		user, err := uc.GetUserById(ctx, userID)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
	})

	// Ошибка репозитория
	t.Run("Repository Error", func(t *testing.T) {
		userID := "validUserID"
		mockAuthRepo.GetUserByIdFunc = func(ctx context.Context, id string) (*domain.User, error) {
			return nil, errors.New("repository error")
		}

		user, err := uc.GetUserById(ctx, userID)
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "repository error")
	})
}
