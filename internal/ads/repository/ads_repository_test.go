package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	ntype "2024_2_FIGHT-CLUB/internal/service/type"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}))
	assert.NoError(t, err)

	return gormDB, mock
}

func TestGetAllPlaces(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock := SetupMockDB(t)
	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	filter := domain.AdFilter{
		Location: "5km",
		Rating:   "4",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "ads"."id","ads"."location_main","ads"."location_street","ads"."position","ads"."images","ads"."author_uuid","ads"."publication_date","ads"."available_dates","ads"."distance" FROM "ads" JOIN users ON ads.author_uuid = users.uuid WHERE distance <= $1 AND users.score >= $2`)).
		WithArgs(5, int64(4)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "location_main", "location_street", "position", "images", "author_uuid", "publication_date", "available_dates", "distance"}).
			AddRow("1", "New York", "5th Avenue", append(make(ntype.Float64Array, 0), 40.7128, -74.0060), "[]", "123", "2024-10-20", "[]", 4.5).
			AddRow("2", "Los Angeles", "Sunset Blvd", append(make(ntype.Float64Array, 0), 34.0522, -118.2437), "[]", "456", "2024-10-18", "[]", 3.0))

	ads, err := adRepo.GetAllPlaces(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(ads))
	assert.Equal(t, "New York", ads[0].LocationMain)
	assert.Equal(t, "Los Angeles", ads[1].LocationMain)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlaceById(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock := SetupMockDB(t)
	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "1"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE id = $1 ORDER BY "ads"."id" LIMIT $2`)).
		WithArgs(adId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "location_main", "location_street", "position", "images", "author_uuid", "publication_date", "available_dates", "distance"}).
			AddRow("1", "New York", "5th Avenue", append(make(ntype.Float64Array, 0), 40.7128, -74.0060), "[]", "123", "2024-10-20", "[]", 4.5))

	ad, err := adRepo.GetPlaceById(ctx, adId)

	assert.NoError(t, err)
	assert.Equal(t, "New York", ad.LocationMain)
	assert.Equal(t, "123", ad.AuthorUUID)
	assert.Equal(t, "5th Avenue", ad.LocationStreet)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := SetupMockDB(t)
	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "1"
	userId := "123"
	ad := &domain.Ad{
		ID:           "1",
		LocationMain: "Los Angeles",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE id = $1 ORDER BY "ads"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "location_main", "author_uuid"}).
			AddRow("1", "New York", "123"))

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "id"=$1,"location_main"=$2 WHERE "id" = $3`)).
		WithArgs("1", "Los Angeles", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := adRepo.UpdatePlace(ctx, ad, adId, userId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock := SetupMockDB(t)
	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "1"
	userId := "123"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE id = $1 ORDER BY "ads"."id" LIMIT $2`)).
		WithArgs(adId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "location_main", "author_uuid"}).
			AddRow("1", "New York", "123"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ads" WHERE "ads"."id" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := adRepo.DeletePlace(ctx, adId, userId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock := SetupMockDB(t)
	defer func() {
		db, _ := db.DB()
		db.Close()
	}()

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	ad := &domain.Ad{
		ID:              "1",
		LocationMain:    "New York",
		LocationStreet:  "5th Avenue",
		Position:        append(make(ntype.Float64Array, 0), 40.7128, -74.0060),
		Images:          append(make(ntype.StringArray, 0), "images1.jpg", "images2.jpg"),
		AuthorUUID:      "123",
		PublicationDate: "2024-10-20",
		AvailableDates:  append(make(ntype.StringArray, 0), "2023-10-11", "2023-10-20"),
		Distance:        4.5,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "ads" ("id","location_main","location_street","position","images","author_uuid","publication_date","available_dates","distance") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs("1", "New York", "5th Avenue", append(make(ntype.Float64Array, 0), 40.7128, -74.0060), append(make(ntype.StringArray, 0), "images1.jpg", "images2.jpg"), "123", "2024-10-20", append(make(ntype.StringArray, 0), "2023-10-11", "2023-10-20"), 4.5).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := adRepo.CreatePlace(ctx, ad)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
