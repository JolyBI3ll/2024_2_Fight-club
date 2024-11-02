package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
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
