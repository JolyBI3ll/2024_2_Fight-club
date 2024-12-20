package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	ntype "2024_2_FIGHT-CLUB/internal/service/type"
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

func setupDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	return gormDB, mock, err
}

func TestGetAllPlaces(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	require.NoError(t, err)

	repo := NewAdRepository(db)

	filter := domain.AdFilter{
		Location:    "",
		Rating:      "",
		NewThisWeek: "",
		HostGender:  "",
		GuestCount:  "",
	}

	userId := "12345"

	// Фиксированная дата для теста
	fixedDate := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	query := `
		SELECT ads.*, cities.title as "CityName", ad_available_dates."availableDateFrom" as "AdDateFrom", ad_available_dates."availableDateTo" as "AdDateTo"
		FROM "ads"
		JOIN cities ON ads."cityId" = cities.id
		JOIN users ON ads."authorUUID" = users.uuid
		JOIN ad_available_dates ON ad_available_dates."adId" = ads.uuid
	`

	adRows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address", "publicationDate", "description", "roomsNumber", "viewsCount",
		"CityName", "availableDateFrom", "availableDateTo",
	}).AddRow("some-uuid", 1, "author-uuid", "Some Address", fixedDate, "Some Description", 3, 0, "City Name", fixedDate, fixedDate)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(adRows)

	favoritesQuery := `SELECT "adId" FROM "favorites" WHERE "userId" = $1`
	favoritesRows := sqlmock.NewRows([]string{"adId"}).AddRow("id1").AddRow("id2")
	mock.ExpectQuery(regexp.QuoteMeta(favoritesQuery)).WillReturnRows(favoritesRows)

	imagesQuery := `SELECT * FROM "images" WHERE "adId" = $1`
	imageRows := sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).
		AddRow(1, "some-uuid", "images/image1.jpg").
		AddRow(2, "some-uuid", "images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).WithArgs("some-uuid").WillReturnRows(imageRows)

	userQuery := `SELECT * FROM "users" WHERE uuid = $1`
	userRows := sqlmock.NewRows([]string{
		"uuid", "username", "password", "email", "name", "score", "avatar", "sex", "guestCount", "birthdate",
	}).AddRow("author-uuid", "test_username", "some_password", "test@example.com", "Test User", 4.5, "avatar_url", "M", 2, fixedDate)
	mock.ExpectQuery(regexp.QuoteMeta(userQuery)).WithArgs("author-uuid").WillReturnRows(userRows)

	adQuery := `SELECT * FROM "ad_rooms" WHERE "adId" = $1`
	addRows := sqlmock.NewRows([]string{"id", "adId", "type", "squaremeters"}).AddRow(1, "id1", "some-type", 12)
	mock.ExpectQuery(regexp.QuoteMeta(adQuery)).WithArgs("some-uuid").WillReturnRows(addRows)

	ads, err := repo.GetAllPlaces(context.Background(), filter, userId)

	require.NoError(t, err)
	assert.Len(t, ads, 1)
	assert.Equal(t, "some-uuid", ads[0].UUID)
	assert.Equal(t, "Some Address", ads[0].Address)
	assert.Equal(t, "City Name", ads[0].CityName)
	assert.Equal(t, 3, ads[0].RoomsNumber)

	assert.ElementsMatch(t, []domain.ImageResponse{
		{
			ID:        1,
			ImagePath: "images/image1.jpg",
		},
		{
			ID:        2,
			ImagePath: "images/image2.jpg",
		},
	}, ads[0].Images)

	// Фиксированная дата birthdate для проверки
	assert.Equal(t, domain.UserResponce{
		Name:       "Test User",
		Avatar:     "avatar_url",
		Rating:     4.5,
		GuestCount: 2,
		Sex:        "M",
		Birthdate:  fixedDate,
	}, ads[0].AdAuthor)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetAllPlaces_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	userId := "12345"

	repo := NewAdRepository(db)
	filter := domain.AdFilter{}
	query := `
		SELECT ads.*, cities.title as "CityName", ad_available_dates."availableDateFrom" as "AdDateFrom", ad_available_dates."availableDateTo" as "AdDateTo"
		FROM "ads"
		JOIN cities ON ads."cityId" = cities.id
		JOIN users ON ads."authorUUID" = users.uuid
		JOIN ad_available_dates ON ad_available_dates."adId" = ads.uuid
	`
	mock.ExpectQuery(query).
		WillReturnError(errors.New("db error"))

	ads, err := repo.GetAllPlaces(context.Background(), filter, userId)
	assert.Error(t, err)
	assert.Nil(t, ads)
}

func TestGetPlaceById(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	expectedAd := domain.GetAllAdsResponse{
		UUID:       "some-uuid",
		CityID:     1,
		AuthorUUID: "author-uuid",
		Address:    "Some Address",
	}

	// Step 2: Define Mock Database Expectations
	fixedDate := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	query := "SELECT ads.*, cities.title as \"CityName\", ad_available_dates.\"availableDateFrom\" as \"AdDateFrom\", ad_available_dates.\"availableDateTo\" as \"AdDateTo\" FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id JOIN ad_available_dates ON ad_available_dates.\"adId\" = ads.uuid WHERE \"adId\" = $1"

	adRows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address", "publicationDate", "description", "roomsNumber", "viewsCount",
		"CityName", "availableDateFrom", "availableDateTo",
	}).AddRow("some-uuid", 1, "author-uuid", "Some Address", fixedDate, "Some Description", 3, 0, "City Name", fixedDate, fixedDate)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("some-uuid").WillReturnRows(adRows)

	imagesQuery := `SELECT * FROM "images" WHERE "adId" = $1`
	imageRows := sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).
		AddRow(1, "some-uuid", "images/image1.jpg").
		AddRow(2, "some-uuid", "images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).WithArgs("some-uuid").WillReturnRows(imageRows)

	userQuery := "SELECT * FROM \"users\" WHERE uuid = $1"
	userRows := sqlmock.NewRows([]string{
		"uuid", "username", "password", "email", "name", "score", "avatar", "sex", "guestCount", "birthdate",
	}).AddRow("author-uuid", "test_username", "some_password", "test@example.com", "Test User", 4.5, "avatar_url", "M", 2, fixedDate)
	mock.ExpectQuery(regexp.QuoteMeta(userQuery)).WithArgs("author-uuid").WillReturnRows(userRows)

	adQuery := `SELECT * FROM "ad_rooms" WHERE "adId" = $1`
	addRows := sqlmock.NewRows([]string{"id", "adId", "type", "squaremeters"}).AddRow(1, "id1", "some-type", 12)
	mock.ExpectQuery(regexp.QuoteMeta(adQuery)).WithArgs("some-uuid").WillReturnRows(addRows)

	ad, err := repo.GetPlaceById(context.Background(), "some-uuid")

	require.NoError(t, err)
	assert.Equal(t, expectedAd.UUID, ad.UUID)
	assert.Equal(t, expectedAd.CityID, ad.CityID)
	assert.Equal(t, expectedAd.AuthorUUID, ad.AuthorUUID)
	assert.Equal(t, expectedAd.Address, ad.Address)
	assert.ElementsMatch(t, []domain.ImageResponse{
		{
			ID:        1,
			ImagePath: "images/image1.jpg",
		},
		{
			ID:        2,
			ImagePath: "images/image2.jpg",
		},
	}, ad.Images)

	assert.Equal(t, domain.UserResponce{
		Name:       "Test User",
		Avatar:     "avatar_url",
		Rating:     4.5,
		GuestCount: 2,
		Sex:        "M",
		Birthdate:  fixedDate,
	}, ad.AdAuthor)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetPlaceById_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT ads.*, cities.title as \"CityName\", ad_available_dates.\"availableDateFrom\" as \"AdDateFrom\", ad_available_dates.\"availableDateTo\" as \"AdDateTo\" FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id JOIN ad_available_dates ON ad_available_dates.\"adId\" = ads.uuid WHERE \"adId\" = $1")).
		WithArgs("ad1").
		WillReturnError(errors.New("db error"))

	ad, err := repo.GetPlaceById(context.Background(), "ad1")
	assert.Error(t, err)
	assert.Empty(t, ad)
}

func TestAdRepository_UpdateViewsCount_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)
	defer func(db *gorm.DB) {
		_, err := db.DB()
		if err != nil {
			return
		}
	}(db)

	repo := NewAdRepository(db)
	ctx := context.Background()

	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-request-id")

	ad := domain.GetAllAdsResponse{
		UUID:       "ad-uuid-123",
		ViewsCount: 5,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "uuid"=$1,"viewsCount"=$2 WHERE uuid = $3`)).
		WithArgs(ad.UUID, ad.ViewsCount+1, ad.UUID). // Увеличение на 1
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updatedAd, err := repo.UpdateViewsCount(ctx, ad)

	assert.NoError(t, err)
	assert.Equal(t, ad.ViewsCount+1, updatedAd.ViewsCount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdRepository_UpdateViewsCount_DBError(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)
	defer func(db *gorm.DB) {
		_, err := db.DB()
		if err != nil {
			return
		}
	}(db)

	repo := NewAdRepository(db)
	ctx := context.Background()

	ctx = context.WithValue(ctx, middleware.RequestIDKey, "test-request-id")

	ad := domain.GetAllAdsResponse{
		UUID:       "ad-uuid-123",
		ViewsCount: 5,
	}
	updatedViewsCountFail := 5

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "uuid"=$1,"viewsCount"=$2 WHERE uuid = $3`)).
		WithArgs(ad.UUID, ad.ViewsCount+1, ad.UUID).
		WillReturnError(errors.New("error updating views count"))
	mock.ExpectRollback()

	_, err = repo.UpdateViewsCount(ctx, ad)

	assert.Error(t, err)
	assert.Equal(t, "error updating views count", err.Error())
	assert.Equal(t, ad.ViewsCount, updatedViewsCountFail)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdRepository_UpdateFavoritesCount_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	favoritesCount := int64(10)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "favorites" WHERE "adId" = $1`)).
		WithArgs(adID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(favoritesCount))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "likesCount"=$1 WHERE uuid = $2`)).
		WithArgs(favoritesCount, adID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdateFavoritesCount(ctx, adID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdRepository_UpdateFavoritesCount_CountError(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "favorites" WHERE "adId" = $1`)).
		WithArgs(adID).
		WillReturnError(errors.New("error counting favorites"))

	err = repo.UpdateFavoritesCount(ctx, adID)

	assert.Error(t, err)
	assert.Equal(t, "error counting favorites", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdRepository_UpdateFavoritesCount_UpdateError(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	favoritesCount := int64(10)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "favorites" WHERE "adId" = $1`)).
		WithArgs(adID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(favoritesCount))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "likesCount"=$1 WHERE uuid = $2`)).
		WithArgs(favoritesCount, adID).
		WillReturnError(errors.New("error updating favorites count"))
	mock.ExpectRollback()

	err = repo.UpdateFavoritesCount(ctx, adID)

	assert.Error(t, err)
	assert.Equal(t, "error updating favorites count", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "existing-ad-id"
	userId := "author-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "Новый город",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
		DateFrom:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		DateTo:      time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		HasBalcony:  true,
		HasGas:      true,
		HasElevator: true,
	}

	adRows := sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}).
		AddRow("existing-ad-id", "author-id", 1, "Old Address", 2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(adRows)

	dateRows := sqlmock.NewRows([]string{"id", "adId", "availableDateFrom", "availableDateTo"}).
		AddRow(1, "existing-ad-id", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ad_available_dates" WHERE "adId" = $1 ORDER BY "ad_available_dates"."id" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(dateRows)

	cityRows := sqlmock.NewRows([]string{"id", "title", "enTitle"}).AddRow(1, "Город", "City").AddRow(2, "Новый Город", "New city")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE title = $1 ORDER BY "cities"."id" LIMIT $2`)).
		WithArgs("Новый город", 1).WillReturnRows(cityRows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "cityId"=$1,"address"=$2,"description"=$3,"roomsNumber"=$4,"hasBalcony"=$5,"hasElevator"=$6,"hasGas"=$7 WHERE "uuid" = $8`)).
		WithArgs(1, "New Address", "Updated Description", 3, true, true, true, adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "hasBalcony"=$1 WHERE "uuid" = $2`)).
		WithArgs(true, adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "hasElevator"=$1 WHERE "uuid" = $2`)).
		WithArgs(true, adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "hasGas"=$1 WHERE "uuid" = $2`)).
		WithArgs(true, adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ad_available_dates" SET "id"=$1,"adId"=$2,"availableDateFrom"=$3,"availableDateTo"=$4 WHERE "id" = $5`)).
		WithArgs(1, adId, updatedRequest.DateFrom, updatedRequest.DateTo, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ad_rooms" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePlace_AdNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "non-existing-ad-id"
	userId := "author-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "New City",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}))

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)

	assert.Error(t, err)
	assert.Equal(t, errors.New("ad not found"), err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePlace_AdDateNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "existing-ad-id"
	userId := "unauthorized-user-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "New City",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	adRows := sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}).
		AddRow("existing-ad-id", "different-author-id", 1, "Old Address", 2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(adRows)

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)

	assert.Error(t, err)
	assert.Equal(t, errors.New("ad date not found"), err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePlace_UserNotAuthorized(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "existing-ad-id"
	userId := "unauthorized-user-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "New City",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	adRows := sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}).
		AddRow("existing-ad-id", "different-author-id", 1, "Old Address", 2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(adRows)

	dateRows := sqlmock.NewRows([]string{"id", "adId", "availableDateFrom", "availableDateTo"}).
		AddRow(1, "existing-ad-id", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ad_available_dates" WHERE "adId" = $1 ORDER BY "ad_available_dates"."id" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(dateRows)

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)
	assert.Error(t, err)
	assert.Equal(t, errors.New("not owner of ad"), err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeletePlace_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "1"
	userId := "not-owner"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "location_main", "author_uuid"}).
			AddRow("1", "New York", "123"))

	err = adRepo.DeletePlace(ctx, adId, userId)

	assert.Error(t, err)
	assert.Equal(t, "not owner of ad", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	require.NoError(t, err)

	repo := NewAdRepository(db)

	ad := &domain.Ad{
		UUID:         "some-ad-uuid",
		Description:  "A great place",
		Address:      "123 Main Street",
		RoomsNumber:  2,
		SquareMeters: 15,
		Floor:        2,
		BuildingType: "house",
		HasElevator:  true,
		HasGas:       true,
		HasBalcony:   true,
		LikesCount:   5,
		Priority:     10,
	}

	newAd := domain.CreateAdRequest{
		CityName:    "Test City",
		Description: "A great place",
		Address:     "123 Main Street",
		RoomsNumber: 2,
		DateFrom:    time.Now(),
		DateTo:      time.Now().AddDate(0, 0, 7),
		Rooms: []domain.AdRoomsResponse{
			{
				Type:         "room",
				SquareMeters: 14,
			},
		},
		SquareMeters: 15,
		Floor:        2,
		BuildingType: "house",
		HasElevator:  true,
		HasGas:       true,
		HasBalcony:   true,
	}

	user := domain.User{
		UUID:   "user-uuid",
		IsHost: true,
	}

	city := domain.City{
		ID:    1,
		Title: "Test City",
	}

	date := domain.AdAvailableDate{
		ID:                1,
		AdID:              "some-ad-uuid",
		AvailableDateFrom: newAd.DateFrom,
		AvailableDateTo:   newAd.DateTo,
	}

	room := domain.AdRooms{
		ID:           1,
		AdID:         "some-ad-uuid",
		Type:         "room",
		SquareMeters: 14,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
		WithArgs(user.UUID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "isHost"}).
			AddRow(user.UUID, user.IsHost))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE title = $1 ORDER BY "cities"."id" LIMIT $2`)).
		WithArgs(newAd.CityName, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
			AddRow(city.ID, city.Title))

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ads" ("cityId","authorUUID","address","publicationDate","description","roomsNumber","viewsCount","squareMeters","floor","buildingType","hasBalcony","hasElevator","hasGas","likesCount","priority","endBoostDate","uuid") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17) RETURNING "uuid"`)).
		WithArgs(
			city.ID,          // cityId
			user.UUID,        // authorUUID
			ad.Address,       // address
			sqlmock.AnyArg(), // publicationDate
			ad.Description,   // description
			ad.RoomsNumber,   // roomsNumber
			0,                // viewsCount
			ad.SquareMeters,
			ad.Floor,
			ad.BuildingType,
			ad.HasBalcony,
			ad.HasElevator,
			ad.HasGas,
			ad.LikesCount,
			ad.Priority,
			sqlmock.AnyArg(),
			ad.UUID, // uuid
		).
		WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(ad.UUID))
	mock.ExpectCommit()

	// Ожидание для вставки в таблицу "ad_available_dates"
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ad_available_dates" ("adId","availableDateFrom","availableDateTo") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(
			date.AdID,              // adID
			date.AvailableDateFrom, // availableDateFrom
			date.AvailableDateTo,   // availableDateTo
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(date.ID))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ad_rooms" ("adId","type","squareMeters") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(
			room.AdID,
			room.Type,
			room.SquareMeters,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(date.ID))
	mock.ExpectCommit()

	err = repo.CreatePlace(context.Background(), ad, newAd, user.UUID)
	require.NoError(t, err)

	assert.Equal(t, city.ID, ad.CityID)
	assert.Equal(t, user.UUID, ad.AuthorUUID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCreatePlace_CityNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	newAdReq := domain.CreateAdRequest{
		CityName:    "Unknown City",
		Address:     "123 Main St",
		Description: "A lovely place",
		RoomsNumber: 3,
		DateFrom:    time.Now(),
		DateTo:      time.Now().AddDate(0, 0, 7),
	}

	ad := &domain.Ad{
		UUID:            "ad-uuid-123",
		AuthorUUID:      "user-uuid-456",
		Address:         "123 Main St",
		Description:     "A lovely place",
		RoomsNumber:     3,
		PublicationDate: time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
		WithArgs(ad.AuthorUUID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "isHost"}).
			AddRow(ad.AuthorUUID, true))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE title = $1 ORDER BY "cities"."id" LIMIT $2`)).
		WithArgs(newAdReq.CityName, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	err = adRepo.CreatePlace(ctx, ad, newAdReq, ad.AuthorUUID)

	if err == nil || err.Error() != "error finding city" {
		t.Errorf("Expected 'error finding city', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestCreatePlace_UserNotHost(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	newAdReq := domain.CreateAdRequest{
		CityName:    "Test City",
		Address:     "123 Main St",
		Description: "A lovely place",
		RoomsNumber: 3,
		DateFrom:    time.Now(),
		DateTo:      time.Now().AddDate(0, 0, 7),
	}

	ad := &domain.Ad{
		UUID:            "ad-uuid-123",
		AuthorUUID:      "user-uuid-456",
		Address:         "123 Main St",
		Description:     "A lovely place",
		RoomsNumber:     3,
		PublicationDate: time.Now(),
	}

	// Mock query to find the user - user is not a host
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE uuid = $1 ORDER BY "users"."uuid" LIMIT $2`)).
		WithArgs(ad.AuthorUUID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "isHost"}).
			AddRow(ad.AuthorUUID, false))

	err = adRepo.CreatePlace(ctx, ad, newAdReq, ad.AuthorUUID)

	if err == nil || err.Error() != "user is not host" {
		t.Errorf("Expected 'user is not host', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetPlacesPerCity_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	city := "New York"

	rows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address",
		"publicationDate", "description", "roomsNumber",
		"avatar", "name", "rating", "cityName",
	}).
		AddRow("ad-uuid-123", 1, "author-uuid-456", "123 Main St",
			time.Now(), "A lovely place", 3, "avatar.png", "John Doe", 4.5, "New York")

	// Mock select ads
	mock.ExpectQuery(regexp.QuoteMeta("SELECT ads.*, cities.title as \"CityName\" FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id WHERE cities.\"enTitle\" = $1")).
		WithArgs(city).
		WillReturnRows(rows)

	imagesQuery := "SELECT * FROM \"images\" WHERE \"adId\" = $1"
	imageRows := sqlmock.NewRows(ntype.StringArray{"imageUrl"}).AddRow("images/image1.jpg").AddRow("images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).WithArgs("ad-uuid-123").WillReturnRows(imageRows)

	query2 := "SELECT * FROM \"users\" WHERE uuid = $1"
	rows2 := sqlmock.NewRows([]string{"uuid", "username", "password", "email", "username"}).
		AddRow("some-uuid", "test_username", "some_password", "test@example.com", "test_username")
	mock.ExpectQuery(regexp.QuoteMeta(query2)).WillReturnRows(rows2)

	adQuery := "SELECT * FROM \"ad_rooms\" WHERE \"adId\" = $1"
	addRows := sqlmock.NewRows([]string{"id", "adId", "type", "squaremeters"}).AddRow(1, "id1", "some-type", 12)
	mock.ExpectQuery(regexp.QuoteMeta(adQuery)).WithArgs("ad-uuid-123").WillReturnRows(addRows)

	ads, err := adRepo.GetPlacesPerCity(ctx, city)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ads) != 1 {
		t.Errorf("Expected 1 ad, got %d", len(ads))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetPlacesPerCity_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	city := "Unknown City"

	// Mock select ads with error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT ads.*, cities.title as \"CityName\" FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id WHERE cities.\"enTitle\" = $1")).
		WithArgs(city).
		WillReturnError(gorm.ErrInvalidData)

	ads, err := adRepo.GetPlacesPerCity(ctx, city)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if ads != nil {
		t.Errorf("Expected nil ads, got %v", ads)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSaveImages_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adUUID := "ad-uuid-123"
	imagePaths := []string{"img1.png", "img2.png"}

	for _, path := range imagePaths {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("adId","imageUrl") VALUES ($1,$2) RETURNING "id"`)).
			WithArgs(adUUID, path).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
	}

	err = adRepo.SaveImages(ctx, adUUID, imagePaths)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSaveImages_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adUUID := "ad-uuid-123"
	imagePaths := []string{"img1.png", "img2.png"}

	// Mock insert first image
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("adId","imageUrl") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(adUUID, imagePaths[0]).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	// Mock insert second image with error
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("adId","imageUrl") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(adUUID, imagePaths[1]).
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err = adRepo.SaveImages(ctx, adUUID, imagePaths)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetAdImages_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "ad-uuid-123"

	rows := sqlmock.NewRows([]string{"imageUrl"}).
		AddRow("img1.png").
		AddRow("img2.png")

	// Mock select images
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "imageUrl" FROM "images" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnRows(rows)

	images, err := adRepo.GetAdImages(ctx, adId)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(images))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetAdImages_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "ad-uuid-123"

	// Mock select images with error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "imageUrl" FROM "images" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnError(gorm.ErrInvalidData)

	images, err := adRepo.GetAdImages(ctx, adId)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if images != nil {
		t.Errorf("Expected nil images, got %v", images)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetUserPlaces_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	userId := "author-uuid-456"

	rows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address",
		"publicationDate", "description", "roomsNumber",
		"avatar", "name", "rating", "сityName",
	}).
		AddRow("ad-uuid-123", 1, userId, "123 Main St",
			time.Now(), "A lovely place", 3, "avatar.png", "John Doe", 4.5, "New York")

	// Mock select ads
	mock.ExpectQuery(`SELECT ads\.\*, users\.avatar, users\.name, users\.score as rating, cities\.title as "CityName" FROM "ads" JOIN users ON ads\."authorUUID" = users\.uuid JOIN cities ON ads\."cityId" = cities\.id WHERE users\.uuid = \$1`).
		WithArgs(userId).
		WillReturnRows(rows)

	// Mock select images for each ad
	imageRows := sqlmock.NewRows([]string{"uuid", "adId", "imageUrl"}).
		AddRow("img-uuid-1", "ad-uuid-123", "img1.png").
		AddRow("img-uuid-2", "ad-uuid-123", "img2.png")

	mock.ExpectQuery(`SELECT \* FROM "images" WHERE "adId" = \$1`).
		WithArgs("ad-uuid-123").
		WillReturnRows(imageRows)

	adQuery := "SELECT * FROM \"ad_rooms\" WHERE \"adId\" = $1"
	addRows := sqlmock.NewRows([]string{"id", "adId", "type", "squaremeters"}).AddRow(1, "id1", "some-type", 12)
	mock.ExpectQuery(regexp.QuoteMeta(adQuery)).WithArgs("ad-uuid-123").WillReturnRows(addRows)

	ads, err := adRepo.GetUserPlaces(ctx, userId)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ads) != 1 {
		t.Errorf("Expected 1 ad, got %d", len(ads))
	}

	if len(ads[0].Images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(ads[0].Images))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetUserPlaces_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	userId := "author-uuid-456"

	rows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address",
		"publicationDate", "description", "roomsNumber",
		"avatar", "name", "rating", "cityName",
	}).
		AddRow("ad-uuid-123", 1, userId, "123 Main St",
			time.Now(), "A lovely place", 3, "avatar.png", "John Doe", 4.5, "New York")

	// Mock select ads
	mock.ExpectQuery(`SELECT ads\.\*, users\.avatar, users\.name, users\.score as rating, cities\.title as "CityName" FROM "ads" JOIN users ON ads\."authorUUID" = users\.uuid JOIN cities ON ads\."cityId" = cities\.id WHERE users\.uuid = \$1`).
		WithArgs(userId).
		WillReturnRows(rows)

	// Mock select images with error
	mock.ExpectQuery(`SELECT \* FROM "images" WHERE "adId" = \$1`).
		WithArgs("ad-uuid-123").
		WillReturnError(gorm.ErrInvalidData)

	ads, err := adRepo.GetUserPlaces(ctx, userId)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if ads != nil {
		t.Errorf("Expected nil ads, got %v", ads)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestAdRepository_AddToFavorites_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	userID := "user-uuid-456"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(adID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "favorites" ("adId","userId") VALUES ($1,$2)`)).
		WithArgs(adID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.AddToFavorites(ctx, adID, userID)

	assert.NoError(t, err, "AddToFavorites should succeed without error")
	assert.NoError(t, mock.ExpectationsWereMet(), "All database calls should be met")
}

func TestAdRepository_AddToFavorites_AdNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	userID := "user-uuid-456"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	err = repo.AddToFavorites(ctx, adID, userID)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet(), "All database calls should be met")
}

func TestAdRepository_AddToFavorites_ErrorOnCreateFavorite(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	userID := "user-uuid-456"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(adID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "favorites" ("adId","userId") VALUES ($1,$2)`)).
		WithArgs(adID, userID).
		WillReturnError(errors.New("error creating favorite"))
	mock.ExpectRollback()

	err = repo.AddToFavorites(ctx, adID, userID)

	assert.Error(t, err)
	assert.Equal(t, "error create favorite", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet(), "All database calls should be met")
}

func TestAdRepository_DeleteFromFavorites_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	userID := "user-uuid-456"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(adID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "favorites" WHERE ("favorites"."adId","favorites"."userId") IN (($1,$2))`)).
		WithArgs(adID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.DeleteFromFavorites(ctx, adID, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdRepository_DeleteFromFavorites_AdNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	userID := "user-uuid-456"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adID, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	err = repo.DeleteFromFavorites(ctx, adID, userID)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAdRepository_DeleteFromFavorites_DeleteError(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	adID := "ad-uuid-123"
	userID := "user-uuid-456"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(adID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "favorites" WHERE ("favorites"."adId","favorites"."userId") IN (($1,$2))`)).
		WithArgs(adID, userID).
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	err = repo.DeleteFromFavorites(ctx, adID, userID)

	assert.Error(t, err)
	assert.Equal(t, "error create favorite", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserFavorites(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	require.NoError(t, err)

	repo := NewAdRepository(db)
	userId := "user-uuid-123"

	// Фиксированная дата для тестов
	fixedDate := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Основной запрос на выборку избранных объявлений
	query := `
		SELECT ads.*, favorites."userId" AS "FavoriteUserId", cities.title as "CityName", ad_available_dates."availableDateFrom" as "AdDateFrom", ad_available_dates."availableDateTo" as "AdDateTo"
		FROM "ads"
		JOIN favorites ON favorites."adId" = ads.uuid
		JOIN cities ON ads."cityId" = cities.id
		JOIN ad_available_dates ON ad_available_dates."adId" = ads.uuid
		WHERE favorites."userId" = $1
	`

	adRows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address", "publicationDate", "description", "roomsNumber", "viewsCount",
		"FavoriteUserId", "CityName", "AdDateFrom", "AdDateTo",
	}).
		AddRow("ad-uuid-123", 1, "author-uuid-456", "Test Address", fixedDate, "Test Description", 3, 10,
			"user-uuid-123", "CityName", fixedDate, fixedDate)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId).
		WillReturnRows(adRows)

	// Запрос на получение картинок
	imagesQuery := `SELECT * FROM "images" WHERE "adId" = $1`
	imageRows := sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).
		AddRow(1, "ad-uuid-123", "images/image1.jpg").
		AddRow(2, "ad-uuid-123", "images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).
		WithArgs("ad-uuid-123").
		WillReturnRows(imageRows)

	// Запрос на получение данных пользователя
	userQuery := `SELECT * FROM "users" WHERE uuid = $1`
	userRows := sqlmock.NewRows([]string{
		"uuid", "name", "avatar", "score", "sex", "guestCount", "birthdate",
	}).
		AddRow("author-uuid-456", "Test User", "avatar_url", 4.8, "F", 5, fixedDate)
	mock.ExpectQuery(regexp.QuoteMeta(userQuery)).
		WithArgs("author-uuid-456").
		WillReturnRows(userRows)

	// Запрос на получение комнат
	roomsQuery := `SELECT * FROM "ad_rooms" WHERE "adId" = $1`
	roomRows := sqlmock.NewRows([]string{"id", "adId", "type", "squareMeters"}).
		AddRow(1, "ad-uuid-123", "Bedroom", 25).
		AddRow(2, "ad-uuid-123", "Living Room", 40)
	mock.ExpectQuery(regexp.QuoteMeta(roomsQuery)).
		WithArgs("ad-uuid-123").
		WillReturnRows(roomRows)

	// Вызов метода
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	ads, err := repo.GetUserFavorites(ctx, userId)

	// Проверка результата
	require.NoError(t, err)
	assert.Len(t, ads, 1)
	assert.Equal(t, "ad-uuid-123", ads[0].UUID)
	assert.Equal(t, "Test Address", ads[0].Address)
	assert.Equal(t, "CityName", ads[0].CityName)
	assert.Equal(t, fixedDate, ads[0].AdDateFrom)
	assert.Equal(t, 3, ads[0].RoomsNumber)

	assert.ElementsMatch(t, []domain.ImageResponse{
		{ID: 1, ImagePath: "images/image1.jpg"},
		{ID: 2, ImagePath: "images/image2.jpg"},
	}, ads[0].Images)

	assert.Equal(t, domain.UserResponce{
		Name:       "Test User",
		Avatar:     "avatar_url",
		Rating:     4.8,
		GuestCount: 5,
		Sex:        "F",
		Birthdate:  fixedDate,
	}, ads[0].AdAuthor)

	assert.ElementsMatch(t, []domain.AdRoomsResponse{
		{Type: "Bedroom", SquareMeters: 25},
		{Type: "Living Room", SquareMeters: 40},
	}, ads[0].Rooms)

	// Убедиться, что все ожидания выполнены
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserFavoritesFail(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	userId := "12345"

	repo := NewAdRepository(db)
	query := `
		SELECT ads.*, favorites."userId" AS "FavoriteUserId", cities.title as "CityName", ad_available_dates."availableDateFrom" as "AdDateFrom", ad_available_dates."availableDateTo" as "AdDateTo"
		FROM "ads"
		JOIN favorites ON favorites."adId" = ads.uuid
		JOIN cities ON ads."cityId" = cities.id
		JOIN ad_available_dates ON ad_available_dates."adId" = ads.uuid
		WHERE favorites."userId" = $1
	`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userId).
		WillReturnError(errors.New("db error"))

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	ads, err := repo.GetUserFavorites(ctx, userId)
	assert.Error(t, err)
	assert.Nil(t, ads)
}

func TestUpdatePriority(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	require.NoError(t, err)

	repo := NewAdRepository(db)

	adId := "ad-uuid-123"
	userId := "user-uuid-456"
	amount := 5
	now := time.Now()

	selectAdQuery := `SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`
	adRow := sqlmock.NewRows([]string{"uuid", "authorUUID", "priority", "endBoostDate"}).
		AddRow(adId, userId, 1, now)

	mock.ExpectQuery(regexp.QuoteMeta(selectAdQuery)).
		WithArgs(adId, 1).
		WillReturnRows(adRow)

	mock.ExpectBegin()
	updatePriorityQuery := `UPDATE "ads" SET "priority"=$1 WHERE "uuid" = $2`
	mock.ExpectExec(regexp.QuoteMeta(updatePriorityQuery)).
		WithArgs(amount, adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	updateBoostDateQuery := `UPDATE "ads" SET "endBoostDate"=$1 WHERE "uuid" = $2`
	mock.ExpectExec(regexp.QuoteMeta(updateBoostDateQuery)).
		WithArgs(sqlmock.AnyArg(), adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	err = repo.UpdatePriority(ctx, adId, userId, amount)

	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestResetExpiredPriorities(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	require.NoError(t, err)

	repo := NewAdRepository(db)

	updateQuery := `
		UPDATE "ads" SET "endBoostDate"=$1,"priority"=$2 WHERE "endBoostDate" <= $3
	`

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WithArgs(nil, 0, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")
	err = repo.ResetExpiredPriorities(ctx)

	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteAdImage(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer func() {
		err := logger.SyncLoggers()
		if err != nil {
			return
		}
	}()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	// Тест успешного удаления изображения
	t.Run("success", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"
		imageUrl := "/images/image.jpg"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"images\" WHERE id = $1 AND \"adId\" = $2 ORDER BY \"images\".\"id\" LIMIT $3")).
			WithArgs(imageId, adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).AddRow(imageId, adId, imageUrl))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM \"images\"").
			WithArgs(imageId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.NoError(t, err)
		assert.Equal(t, imageUrl, result)

		// Проверяем все ожидания выполнены
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Тест случая, когда объявление не найдено
	t.Run("ad not found", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "ad not found", err.Error())
	})

	// Тест случая, когда пользователь не является владельцем объявления
	t.Run("not owner of ad", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"
		wrongUUID := "another-user-uuid"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, wrongUUID))

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "not owner of ad", err.Error())
	})

	// Тест случая, когда изображение не найдено
	t.Run("image not found", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"images\" WHERE id = $1 AND \"adId\" = $2 ORDER BY \"images\".\"id\" LIMIT $3")).
			WithArgs(imageId, adId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "image not found", err.Error())
	})

	// Тест случая ошибки при удалении изображения
	t.Run("error deleting image", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"
		imageUrl := "/images/image.jpg"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"images\" WHERE id = $1 AND \"adId\" = $2 ORDER BY \"images\".\"id\" LIMIT $3")).
			WithArgs(imageId, adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).AddRow(imageId, adId, imageUrl))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM \"images\"").
			WithArgs(imageId).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "error deleting image from database", err.Error())
	})
}
