package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"reflect"
	"regexp"
	"testing"
)

func TestGetCitiesSuccess(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was encountered while creating SQLMock database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was encountered while creating GORM DB connection", err)
	}

	cityRepo := NewCityRepository(gormDB)

	rows := sqlmock.NewRows([]string{"id", "title", "enTitle", "description"}).
		AddRow(1, "Москва", "moscow", "A large city in Russia.")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities"`)).WillReturnRows(rows)

	ctx := context.TODO()
	cities, err := cityRepo.GetCities(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedCities := []domain.City{
		{ID: 1, Title: "Москва", EnTitle: "moscow", Description: "A large city in Russia."},
	}
	if len(cities) != 1 || !reflect.DeepEqual(cities[0], expectedCities[0]) {
		t.Errorf("expected %v, got %v", expectedCities, cities)
	}
}

func TestGetCitiesFailure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was encountered while creating SQLMock database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was encountered while creating GORM DB connection", err)
	}

	cityRepo := NewCityRepository(gormDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities"`)).WillReturnError(sql.ErrConnDone)

	ctx := context.TODO()
	_, err = cityRepo.GetCities(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}

	expectedErrorMsg := sql.ErrConnDone.Error()
	if err.Error() != expectedErrorMsg {
		t.Errorf("expected error message %v, got %v", expectedErrorMsg, err.Error())
	}
}

func TestGetCityByEnName(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was encountered while creating SQLMock database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was encountered while creating GORM DB connection", err)
	}

	cityRepo := NewCityRepository(gormDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"cities\" WHERE \"enTitle\" = $1 ORDER BY \"cities\".\"id\" LIMIT $2")).
			WithArgs("New York", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title", "enTitle", "description", "image"}).
				AddRow(1, "Нью-Йорк", "New York", "A large city", "image_url"))

		city, err := cityRepo.GetCityByEnName(ctx, "New York")
		assert.NoError(t, err)
		assert.Equal(t, domain.City{
			ID:          1,
			Title:       "Нью-Йорк",
			EnTitle:     "New York",
			Description: "A large city",
			Image:       "image_url",
		}, city)
	})

	t.Run("City Not Found", func(t *testing.T) {
		// Ожидаем, что не будет результатов
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"cities\" WHERE \"enTitle\" = $1 ORDER BY \"cities\".\"id\" LIMIT $2")).
			WithArgs("Unknown City", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		city, err := cityRepo.GetCityByEnName(ctx, "Unknown City")
		assert.Error(t, err)
		assert.Equal(t, domain.City{}, city)
	})

	t.Run("DB Error", func(t *testing.T) {
		// Ожидаем, что будет ошибка базы данных
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"cities\" WHERE \"enTitle\" = $1 ORDER BY \"cities\".\"id\" LIMIT $2")).
			WithArgs("Database Error City", 1).
			WillReturnError(errors.New("database error"))

		city, err := cityRepo.GetCityByEnName(ctx, "Database Error City")
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, domain.City{}, city)
	})

	// Убедимся, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
