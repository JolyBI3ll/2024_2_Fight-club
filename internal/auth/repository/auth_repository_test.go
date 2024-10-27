package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
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

//func TestCreateUser(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer logger.SyncLoggers()
//	gormDB, mock := setupTestDB(t)
//	authRepo := NewAuthRepository(gormDB)
//
//	user := &domain.User{
//		Username:   "test_user",
//		Password:   "password123",
//		Email:      "test@example.com",
//		Name:       "",
//		Score:      0,
//		Avatar:     "",
//		Sex:        0,
//		GuestCount: 0,
//		Birthdate:  time.Time{}, // или любое желаемое значение
//		Address:    "",
//		IsHost:     false,
//	}
//
//	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test_request_id")
//
//	t.Run("Success", func(t *testing.T) {
//		mock.ExpectBegin()
//		mock.ExpectExec("INSERT INTO \"users\" (\"username\",\"password\",\"email\",\"name\",\"score\",\"avatar\",\"sex\",\"guest_count\",\"birthdate\",\"address\",\"is_host\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)").
//			WithArgs(
//				user.Username,
//				user.Password,
//				user.Email,
//				user.Name,
//				user.Score,
//				user.Avatar,
//				user.Sex,
//				user.GuestCount,
//				sqlmock.AnyArg(), // Заменим на sqlmock.AnyArg() для даты рождения
//				user.Address,
//				user.IsHost,
//			).
//			WillReturnResult(sqlmock.NewResult(1, 1))
//		mock.ExpectCommit()
//
//		err := authRepo.CreateUser(ctx, user)
//		assert.NoError(t, err) // Проверка на отсутствие ошибки
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//
//	t.Run("ErrorCreatingUser", func(t *testing.T) {
//		mock.ExpectBegin()
//		mock.ExpectExec("INSERT INTO \"users\"").
//			WithArgs(
//				user.Username,
//				user.Password,
//				user.Email,
//				user.Name,
//				user.Score,
//				user.Avatar,
//				user.Sex,
//				user.GuestCount,
//				sqlmock.AnyArg(), // Обязательно указывать sqlmock.AnyArg() для времени
//				user.Address,
//				user.IsHost,
//			).
//			WillReturnError(errors.New("create error"))
//		mock.ExpectRollback()
//
//		err := authRepo.CreateUser(ctx, user)
//		assert.Error(t, err)                         // Проверка на наличие ошибки
//		assert.Equal(t, "create error", err.Error()) // Проверка текста ошибки
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//}

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
		Sex:        'M',
		GuestCount: 5,
		Birthdate:  time.Time{},
		Address:    "Some Address",
		IsHost:     true,
	}

	mock.ExpectBegin()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("username","password","email","name","score","avatar","sex","guest_count","birthdate","address","is_host","uuid") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "uuid"`)).
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
			expectedUser.Address,
			expectedUser.IsHost,
			expectedUser.UUID).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	err := repo.CreateUser(context.Background(), expectedUser)
	assert.Error(t, err)
}

func TestSaveUser_Success(t *testing.T) {
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
		Sex:        'M',
		GuestCount: 5,
		Birthdate:  time.Time{},
		Address:    "Some Address",
		IsHost:     true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"email"=$3,"name"=$4,"score"=$5,"avatar"=$6,"sex"=$7,"guest_count"=$8,"birthdate"=$9,"address"=$10,"is_host"=$11 WHERE "uuid" = $12`)).
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
			expectedUser.Address,
			expectedUser.IsHost,
			expectedUser.UUID,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.SaveUser(context.Background(), expectedUser)
	assert.NoError(t, err)
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
		Sex:        'M',
		GuestCount: 5,
		Birthdate:  time.Time{},
		Address:    "Some Address",
		IsHost:     true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "username"=$1,"password"=$2,"email"=$3,"name"=$4,"score"=$5,"avatar"=$6,"sex"=$7,"guest_count"=$8,"birthdate"=$9,"address"=$10,"is_host"=$11 WHERE "uuid" = $12`)).
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
			expectedUser.Address,
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE UUID = $1 ORDER BY "users"."uuid" LIMIT $2`)).
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
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE UUID = $1 ORDER BY "users"."uuid" LIMIT $2`)).
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
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByName(ctx, username)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "database error", err.Error())
	})

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepository_PutUser(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
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
		Sex:        'M',
		GuestCount: 5,
		Birthdate:  time.Now(),
		Address:    "Some Address",
		IsHost:     true,
	}

	// Тест-кейс 1: Успешное обновление пользователя
	t.Run("Successfully update user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET`).
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
				creds.Address,
				creds.IsHost,
				userID,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.PutUser(ctx, creds, userID)
		assert.NoError(t, err)
	})

	// Тест-кейс 2: Ошибка при обновлении пользователя
	t.Run("Error updating user", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "users" SET`).
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
				creds.Address,
				creds.IsHost,
				userID,
			).WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := repo.PutUser(ctx, creds, userID)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())
	})

	require.NoError(t, mock.ExpectationsWereMet())
}
