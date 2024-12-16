package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"net/http"
)

type MockJwtTokenService struct {
	MockCreate            func(session_id string, tokenExpTime int64) (string, error)
	MockValidate          func(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error)
	MockParseSecretGetter func(token *jwt.Token) (interface{}, error)
}

func (m *MockJwtTokenService) Create(session_id string, tokenExpTime int64) (string, error) {
	return m.MockCreate(session_id, tokenExpTime)
}

func (m *MockJwtTokenService) Validate(tokenString string, expectedSessionId string) (*middleware.JwtCsrfClaims, error) {
	return m.MockValidate(tokenString, expectedSessionId)
}

func (m *MockJwtTokenService) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	return m.MockParseSecretGetter(token)
}

type MockServiceSession struct {
	MockGetUserID      func(ctx context.Context, sessionID string) (string, error)
	MockLogoutSession  func(ctx context.Context, sessionID string) error
	MockCreateSession  func(ctx context.Context, user *domain.User) (string, error)
	MockGetSessionData func(ctx context.Context, sessionID string) (*map[string]interface{}, error)
	MockGetSession     func(r *http.Request) (*sessions.Session, error)
}

func (m *MockServiceSession) GetUserID(ctx context.Context, sessionID string) (string, error) {
	return m.MockGetUserID(ctx, sessionID)
}

func (m *MockServiceSession) LogoutSession(ctx context.Context, sessionID string) error {
	return m.MockLogoutSession(ctx, sessionID)
}

func (m *MockServiceSession) CreateSession(ctx context.Context, user *domain.User) (string, error) {
	return m.MockCreateSession(ctx, user)
}

func (m *MockServiceSession) GetSessionData(ctx context.Context, sessionID string) (*map[string]interface{}, error) {
	return m.MockGetSessionData(ctx, sessionID)
}

func (m *MockServiceSession) GetSession(r *http.Request) (*sessions.Session, error) {
	return m.MockGetSession(r)
}

type MockAdUseCase struct {
	MockGetAllPlaces     func(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error)
	MockGetOnePlace      func(ctx context.Context, adId string, isAuthorized bool) (domain.GetAllAdsResponse, error)
	MockCreatePlace      func(ctx context.Context, place *domain.Ad, fileHeader [][]byte, newPlace domain.CreateAdRequest, userId string) error
	MockUpdatePlace      func(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader [][]byte, updatedPlace domain.UpdateAdRequest) error
	MockDeletePlace      func(ctx context.Context, adId string, userId string) error
	MockGetPlacesPerCity func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error)
	MockGetUserPlaces    func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	MockDeleteAdImage    func(ctx context.Context, adId string, imageId string, userId string) error
}

func (m *MockAdUseCase) DeleteAdImage(ctx context.Context, adId string, imageId string, userId string) error {
	return m.MockDeleteAdImage(ctx, adId, imageId, userId)
}

func (m *MockAdUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetAllPlaces(ctx, filter)
}

func (m *MockAdUseCase) GetOnePlace(ctx context.Context, adId string, isAuthorized bool) (domain.GetAllAdsResponse, error) {
	return m.MockGetOnePlace(ctx, adId, isAuthorized)
}

func (m *MockAdUseCase) CreatePlace(ctx context.Context, place *domain.Ad, fileHeader [][]byte, newPlace domain.CreateAdRequest, userId string) error {
	return m.MockCreatePlace(ctx, place, fileHeader, newPlace, userId)
}

func (m *MockAdUseCase) UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader [][]byte, updatedPlace domain.UpdateAdRequest) error {
	return m.MockUpdatePlace(ctx, place, adId, userId, fileHeader, updatedPlace)
}

func (m *MockAdUseCase) DeletePlace(ctx context.Context, adId string, userId string) error {
	return m.MockDeletePlace(ctx, adId, userId)
}

func (m *MockAdUseCase) GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetPlacesPerCity(ctx, city)
}

func (m *MockAdUseCase) GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetUserPlaces(ctx, userId)
}

type MockAdRepository struct {
	MockGetAllPlaces           func(ctx context.Context, filter domain.AdFilter, userId string) ([]domain.GetAllAdsResponse, error)
	MockGetPlaceById           func(ctx context.Context, adId string) (domain.GetAllAdsResponse, error)
	MockUpdateViewsCount       func(ctx context.Context, ad domain.GetAllAdsResponse) (domain.GetAllAdsResponse, error)
	MockCreatePlace            func(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error
	MockSavePlace              func(ctx context.Context, ad *domain.Ad) error
	MockUpdatePlace            func(ctx context.Context, ad *domain.Ad, adId string, userId string, updatedAd domain.UpdateAdRequest) error
	MockDeletePlace            func(ctx context.Context, adId string, userId string) error
	MockGetPlacesPerCity       func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error)
	MockSaveImages             func(ctx context.Context, adUUID string, imagePaths []string) error
	MockGetAdImages            func(ctx context.Context, adId string) ([]string, error)
	MockGetUserPlaces          func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	MockDeleteAdImage          func(ctx context.Context, adId string, imageId int, userId string) (string, error)
	MockAddToFavorites         func(ctx context.Context, adId string, userId string) error
	MockDeleteFromFavorites    func(ctx context.Context, adId string, userId string) error
	MockGetUserFavorites       func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	MockUpdateFavoritesCount   func(ctx context.Context, adId string) error
	MockUpdatePriority         func(ctx context.Context, adId string, userId string, amount int) error
	MockResetExpiredPriorities func(ctx context.Context) error
}

func (m *MockAdRepository) DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) (string, error) {
	return m.MockDeleteAdImage(ctx, adId, imageId, userId)
}

func (m *MockAdRepository) GetAllPlaces(ctx context.Context, filter domain.AdFilter, userId string) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetAllPlaces(ctx, filter, userId)
}

func (m *MockAdRepository) GetPlaceById(ctx context.Context, adId string) (domain.GetAllAdsResponse, error) {
	return m.MockGetPlaceById(ctx, adId)
}

func (m *MockAdRepository) UpdateViewsCount(ctx context.Context, ad domain.GetAllAdsResponse) (domain.GetAllAdsResponse, error) {
	return m.MockUpdateViewsCount(ctx, ad)
}

func (m *MockAdRepository) CreatePlace(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error {
	return m.MockCreatePlace(ctx, ad, newAd, userId)
}

func (m *MockAdRepository) SavePlace(ctx context.Context, ad *domain.Ad) error {
	return m.MockSavePlace(ctx, ad)
}

func (m *MockAdRepository) UpdatePlace(ctx context.Context, ad *domain.Ad, adId string, userId string, updatedAd domain.UpdateAdRequest) error {
	return m.MockUpdatePlace(ctx, ad, adId, userId, updatedAd)
}

func (m *MockAdRepository) DeletePlace(ctx context.Context, adId string, userId string) error {
	return m.MockDeletePlace(ctx, adId, userId)
}

func (m *MockAdRepository) GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetPlacesPerCity(ctx, city)
}

func (m *MockAdRepository) SaveImages(ctx context.Context, adUUID string, imagePaths []string) error {
	return m.MockSaveImages(ctx, adUUID, imagePaths)
}

func (m *MockAdRepository) GetAdImages(ctx context.Context, adId string) ([]string, error) {
	return m.MockGetAdImages(ctx, adId)
}

func (m *MockAdRepository) GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetUserPlaces(ctx, userId)
}

func (m *MockAdRepository) AddToFavorites(ctx context.Context, adId string, userId string) error {
	return m.MockAddToFavorites(ctx, adId, userId)
}

func (m *MockAdRepository) DeleteFromFavorites(ctx context.Context, adId string, userId string) error {
	return m.MockDeleteFromFavorites(ctx, adId, userId)
}

func (m *MockAdRepository) GetUserFavorites(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetUserFavorites(ctx, userId)
}

func (m *MockAdRepository) UpdateFavoritesCount(ctx context.Context, adId string) error {
	return m.MockUpdateFavoritesCount(ctx, adId)
}

func (m *MockAdRepository) UpdatePriority(ctx context.Context, adId string, userId string, amount int) error {
	return m.MockUpdatePriority(ctx, adId, userId, amount)
}

func (m *MockAdRepository) ResetExpiredPriorities(ctx context.Context) error {
	return m.MockResetExpiredPriorities(ctx)
}

type MockMinioService struct {
	UploadFileFunc func(file []byte, contentType, id string) (string, error)
	DeleteFileFunc func(filePath string) error
}

func (m *MockMinioService) UploadFile(file []byte, contentType, id string) (string, error) {
	return m.UploadFileFunc(file, contentType, id)
}

func (m *MockMinioService) DeleteFile(filePath string) error {
	return m.DeleteFileFunc(filePath)
}
