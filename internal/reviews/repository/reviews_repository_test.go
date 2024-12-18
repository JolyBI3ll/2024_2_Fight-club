package repository

//import (
//	"2024_2_FIGHT-CLUB/domain"
//	"2024_2_FIGHT-CLUB/internal/service/logger"
//	"context"
//	"database/sql"
//	"errors"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/stretchr/testify/assert"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"log"
//	"regexp"
//	"testing"
//)
//
//func TestCreateReview(t *testing.T) {
//	if err := logger.InitLoggers(); err != nil {
//		log.Fatalf("Failed to initialize loggers: %v", err)
//	}
//	defer func() {
//		err := logger.SyncLoggers()
//		if err != nil {
//			return
//		}
//	}()
//	db, mock, err := sqlmock.New()
//	assert.NoError(t, err)
//	defer func(db *sql.DB) {
//		err := db.Close()
//		if err != nil {
//			return
//		}
//	}(db)
//
//	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
//	assert.NoError(t, err)
//
//	repo := NewReviewRepository(gormDB)
//	ctx := context.Background()
//
//	// --- Test Case 1: Host not found ---
//	t.Run("Host not found", func(t *testing.T) {
//		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
//			WithArgs("host123", 1).
//			WillReturnError(gorm.ErrRecordNotFound)
//		review := &domain.Review{UserID: "user123", HostID: "host123", Rating: 4}
//		err := repo.CreateReview(ctx, review)
//
//		assert.EqualError(t, err, "error finding host")
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//
//	// --- Test Case 2: Review already exists ---
//	t.Run("Review already exists", func(t *testing.T) {
//		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
//			WithArgs("user123", 1).
//			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//
//		review := &domain.Review{UserID: "user123", HostID: "user123", Rating: 4}
//		err := repo.CreateReview(ctx, review)
//
//		assert.EqualError(t, err, "review already exist")
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//
//	// --- Test Case 3: Error finding review ---
//	t.Run("Error finding review", func(t *testing.T) {
//		mock.ExpectQuery(`SELECT .* FROM "reviews"`).
//			WithArgs("user123", "host123").
//			WillReturnError(errors.New("some error"))
//
//		review := &domain.Review{UserID: "user123", HostID: "host123", Rating: 4}
//		err := repo.CreateReview(ctx, review)
//
//		assert.EqualError(t, err, "error finding review")
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//
//	// --- Test Case 4: Error creating review ---
//	t.Run("Error creating review", func(t *testing.T) {
//		mock.ExpectQuery(`SELECT .* FROM "users"`).
//			WithArgs("host123").
//			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//
//		mock.ExpectQuery(`SELECT .* FROM "reviews"`).
//			WithArgs("user123", "host123").
//			WillReturnError(gorm.ErrRecordNotFound)
//
//		mock.ExpectExec(`INSERT INTO "reviews"`).
//			WillReturnError(errors.New("insert error"))
//
//		review := &domain.Review{UserID: "user123", HostID: "host123", Rating: 4}
//		err := repo.CreateReview(ctx, review)
//
//		assert.EqualError(t, err, "error creating review")
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//
//	// --- Test Case 5: Error updating host score ---
//	t.Run("Error updating host score", func(t *testing.T) {
//		mock.ExpectQuery(`SELECT .* FROM "users"`).
//			WithArgs("host123").
//			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//
//		mock.ExpectQuery(`SELECT .* FROM "reviews"`).
//			WithArgs("user123", "host123").
//			WillReturnError(gorm.ErrRecordNotFound)
//
//		mock.ExpectExec(`INSERT INTO "reviews"`).
//			WillReturnResult(sqlmock.NewResult(1, 1))
//
//		mock.ExpectExec(`UPDATE "hosts"`). // Эмуляция ошибки при обновлении рейтинга хоста
//							WillReturnError(errors.New("update error"))
//
//		review := &domain.Review{UserID: "user123", HostID: "host123", Rating: 4}
//		err := repo.CreateReview(ctx, review)
//
//		assert.EqualError(t, err, "error updating host score")
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//
//	// --- Test Case 6: Successful creation ---
//	t.Run("Successful creation", func(t *testing.T) {
//		mock.ExpectQuery(`SELECT .* FROM "users"`).
//			WithArgs("host123").
//			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//
//		mock.ExpectQuery(`SELECT .* FROM "reviews"`).
//			WithArgs("user123", "host123").
//			WillReturnError(gorm.ErrRecordNotFound)
//
//		mock.ExpectExec(`INSERT INTO "reviews"`).
//			WillReturnResult(sqlmock.NewResult(1, 1))
//
//		mock.ExpectExec(`UPDATE "hosts"`). // Эмуляция успешного обновления рейтинга
//							WillReturnResult(sqlmock.NewResult(1, 1))
//
//		review := &domain.Review{UserID: "user123", HostID: "host123", Rating: 4}
//		err := repo.CreateReview(ctx, review)
//
//		assert.NoError(t, err)
//		assert.NoError(t, mock.ExpectationsWereMet())
//	})
//}
