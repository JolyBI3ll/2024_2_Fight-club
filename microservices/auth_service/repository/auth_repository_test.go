package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestCreateUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	gormDB, mock := setupTestDB(t)
	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	authRepo := NewAuthRepository(gormDB)

	user := &domain.User{
		UUID:       "test-uuid",
		Username:   "testuser",
		Password:   "password",
		Email:      "test@example.com",
		Name:       "testuser",
		Score:      9.1,
		Avatar:     "images/1.jpg",
		Sex:        "M",
		GuestCount: 5,
		Birthdate:  time.Time{},
		IsHost:     true,
	}

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test_request_id")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password","email","name","score","avatar","sex","guestCount","birthDate","isHost","uuid") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "uuid"`)).
			WithArgs(
				user.Username,
				user.Password,
				user.Email,
				user.Name,
				user.Score,
				user.Avatar,
				user.Sex,
				user.GuestCount,
				user.Birthdate,
				user.IsHost,
				user.UUID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(user.UUID))
		mock.ExpectCommit()

		// Вызов тестируемого метода
		err := authRepo.CreateUser(ctx, user)

		// Проверка результатов
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ErrorCreatingUser", func(t *testing.T) {
		// Настройка ожиданий для ошибки при создании пользователя
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("username","password","email","name","score","avatar","sex","guestCount","birthDate","isHost","uuid") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "uuid"`)).
			WithArgs(
				user.Username,
				user.Password,
				user.Email,
				user.Name,
				user.Score,
				user.Avatar,
				user.Sex,
				user.GuestCount,
				user.Birthdate,
				user.IsHost,
				user.UUID,
			).
			WillReturnError(errors.New("error creating user"))
		mock.ExpectRollback()

		// Вызов тестируемого метода
		err := authRepo.CreateUser(ctx, user)

		// Проверка результатов
		assert.Error(t, err)
		assert.Equal(t, "error creating user", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	expectedUser := &domain.User{
		UUID:       "test-uuid",
		Username:   "testuser",
		Password:   "password",
		Email:      "test@example.com",
		Name:       "testuser",
		Score:      9.1,
		Avatar:     "images/1.jpg",
		Sex:        "M",
		GuestCount: 5,
		Birthdate:  time.Time{},
		IsHost:     true,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" \("username","password","email","name","score","avatar","sex","guestCount","birthDate","isHost","uuid"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9,\$10,\$11\) RETURNING "uuid"`).
		WithArgs(
			expectedUser.Username,
			expectedUser.Password,
			expectedUser.Email,
			expectedUser.Name,
			expectedUser.Score,
			expectedUser.Avatar,
			expectedUser.Sex,
			expectedUser.GuestCount,
			expectedUser.Birthdate,
			expectedUser.IsHost,
			expectedUser.UUID).
		WillReturnError(errors.New("error creating user"))
	mock.ExpectRollback()

	err := repo.CreateUser(context.Background(), expectedUser)

	assert.Error(t, err)
	assert.Equal(t, "error creating user", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestSaveUser_Success(t *testing.T) {
	// Инициализация логгеров
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Настройка тестовой базы данных и sqlmock
	db, mock := setupTestDB(t)
	defer db.DB()

	// Создание репозитория
	repo := NewAuthRepository(db)

	// Определение ожидаемого пользователя
	expectedUser := &domain.User{
		UUID:       "test-uuid",
		Username:   "testuser",
		Password:   "password",
		Email:      "test@example.com",
		Name:       "testuser",
		Score:      9.1,
		Avatar:     "images/1.jpg",
		Sex:        "M",
		GuestCount: 5,
		Birthdate:  time.Time{},
		IsHost:     true,
	}

	// Настройка ожиданий sqlmock
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"email"=$3,"name"=$4,"score"=$5,"avatar"=$6,"sex"=$7,"guestCount"=$8,"birthDate"=$9,"isHost"=$10 WHERE "uuid" = $11`)).
		WithArgs(
			expectedUser.Username,
			expectedUser.Password,
			expectedUser.Email,
			expectedUser.Name,
			expectedUser.Score,
			expectedUser.Avatar,
			expectedUser.Sex,
			expectedUser.GuestCount,
			expectedUser.Birthdate,
			expectedUser.IsHost,
			expectedUser.UUID,
		).WillReturnResult(sqlmock.NewResult(0, 1)) // Для UPDATE LastInsertId = 0
	mock.ExpectCommit()

	// Вызов функции сохранения пользователя
	err := repo.SaveUser(context.Background(), expectedUser)
	assert.NoError(t, err)

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveUser_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	expectedUser := &domain.User{
		UUID:       "test-uuid",
		Username:   "testuser",
		Password:   "password",
		Email:      "test@example.com",
		Name:       "testuser",
		Score:      9.1,
		Avatar:     "images/1.jpg",
		Sex:        "M",
		GuestCount: 5,
		Birthdate:  time.Time{},
		IsHost:     true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"email"=$3,"name"=$4,"score"=$5,"avatar"=$6,"sex"=$7,"guestCount"=$8,"birthDate"=$9,"isHost"=$10 WHERE "uuid" = $11`)).
		WithArgs(
			expectedUser.Username,
			expectedUser.Password,
			expectedUser.Email,
			expectedUser.Name,
			expectedUser.Score,
			expectedUser.Avatar,
			expectedUser.Sex,
			expectedUser.GuestCount,
			expectedUser.Birthdate,
			expectedUser.IsHost,
			expectedUser.UUID,
		).WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	err := repo.SaveUser(context.Background(), expectedUser)
	assert.Error(t, err)
}

func TestGetUserById_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	userID := "test-uuid"
	rows := sqlmock.NewRows([]string{"UUID", "Username", "Password", "Email"}).
		AddRow(userID, "testuser", "password", "test@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
		WithArgs(userID, 1).WillReturnRows(rows)

	user, err := repo.GetUserById(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.UUID)
}

func TestGetUserById_NotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	userID := "test-uuid"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
		WithArgs(userID, 1).WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.GetUserById(context.Background(), userID)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetAllUser_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	rows := sqlmock.NewRows([]string{"UUID", "Username", "Password", "Email"}).
		AddRow("uuid1", "user1", "password1", "user1@example.com").
		AddRow("uuid2", "user2", "password2", "user2@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(rows)

	users, err := repo.GetAllUser(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestGetAllUser_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnError(errors.New("get places error"))

	_, err := repo.GetAllUser(context.Background())
	assert.Error(t, err)
}

func TestAuthRepository_GetUserByName(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	ctx := context.TODO()
	username := "testuser"

	// Тест-кейс 1: Успешное получение пользователя по имени
	t.Run("Successfully get user by name", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"UUID", "username"}).
			AddRow("test-uuid", "testuser")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."uuid" LIMIT $2`)).
			WithArgs(username, 1).
			WillReturnRows(rows)

		user, err := repo.GetUserByName(ctx, username)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
	})

	// Тест-кейс 2: Пользователь не найден
	t.Run("User not found by name", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."uuid" LIMIT $2`)).
			WithArgs(username, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetUserByName(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})

	// Тест-кейс 3: Ошибка при выполнении запроса
	t.Run("Database error when fetching user by name", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."uuid" LIMIT $2`)).
			WithArgs(username, 1).
			WillReturnError(errors.New("error fetching user by name"))

		user, err := repo.GetUserByName(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "error fetching user by name", err.Error())
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepository_PutUser(t *testing.T) {
	// Инициализируем логгеры
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	// Настраиваем тестовую базу данных и sqlmock
	db, mock := setupTestDB(t)
	defer db.DB()

	repo := NewAuthRepository(db)

	ctx := context.TODO()
	userID := "test-uuid"
	creds := &domain.User{
		Username:   "testuser",
		Password:   "password",
		Email:      "test@example.com",
		Name:       "testuser",
		Score:      9.1,
		Avatar:     "images/1.jpg",
		Sex:        "M",
		GuestCount: 5,
		Birthdate:  time.Now(),
		IsHost:     true,
	}

	// Тест-кейс 1: Успешное обновление пользователя
	t.Run("Successfully update user", func(t *testing.T) {
		mock.ExpectBegin()

		// Первый UPDATE запрос
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"email"=$3,"name"=$4,"score"=$5,"avatar"=$6,"sex"=$7,"guestCount"=$8,"birthDate"=$9,"isHost"=$10 WHERE UUID = $11`)).
			WithArgs(
				creds.Username,
				creds.Password,
				creds.Email,
				creds.Name,
				creds.Score,
				creds.Avatar,
				creds.Sex,
				creds.GuestCount,
				creds.Birthdate,
				creds.IsHost,
				userID,
			).WillReturnResult(sqlmock.NewResult(0, 1)) // LastInsertId = 0 для UPDATE
		mock.ExpectCommit()
		// Второй UPDATE запрос
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "isHost"=$1 WHERE UUID = $2`)).
			WithArgs(
				creds.IsHost,
				userID,
			).WillReturnResult(sqlmock.NewResult(0, 1)) // LastInsertId = 0 для UPDATE

		mock.ExpectCommit()

		// Вызов метода сохранения пользователя
		err := repo.PutUser(ctx, creds, userID)
		assert.NoError(t, err)
	})

	// Тест-кейс 2: Ошибка обновления пользователя
	t.Run("Error updating user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"email"=$3,"name"=$4,"score"=$5,"avatar"=$6,"sex"=$7,"guestCount"=$8,"birthDate"=$9,"isHost"=$10 WHERE UUID = $11`)).
			WithArgs(
				creds.Username,
				creds.Password,
				creds.Email,
				creds.Name,
				creds.Score,
				creds.Avatar,
				creds.Sex,
				creds.GuestCount,
				creds.Birthdate,
				creds.IsHost,
				userID,
			).WillReturnError(errors.New("error updating user"))
		mock.ExpectRollback()

		err := repo.PutUser(ctx, creds, userID)

		assert.Error(t, err)
		assert.Equal(t, "error updating user", err.Error())
	})

	// Проверяем, что все ожидания выполнены
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	gormDB, mock := setupTestDB(t)
	defer func() {
		db, _ := gormDB.DB()
		db.Close()
	}()

	authRepo := NewAuthRepository(gormDB)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test_request_id")

	t.Run("Success", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &domain.User{
			UUID:     "test-uuid",
			Username: "testuser",
			Password: "hashedpassword",
			Email:    email,
			Name:     "Test User",
		}

		// Настройка моков
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."uuid" LIMIT $2`)).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "username", "password", "email", "name"}).
				AddRow(expectedUser.UUID, expectedUser.Username, expectedUser.Password, expectedUser.Email, expectedUser.Name))

		// Вызов тестируемого метода
		user, err := authRepo.GetUserByEmail(ctx, email)

		// Проверка результатов
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		email := "notfound@example.com"

		// Настройка моков
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."uuid" LIMIT $2`)).
			WithArgs(email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Вызов тестируемого метода
		user, err := authRepo.GetUserByEmail(ctx, email)

		// Проверка результатов
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		email := "dberror@example.com"

		// Настройка моков
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."uuid" LIMIT $2`)).
			WithArgs(email, 1).
			WillReturnError(errors.New("database error"))

		// Вызов тестируемого метода
		user, err := authRepo.GetUserByEmail(ctx, email)

		// Проверка результатов
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "error fetching user by email", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
