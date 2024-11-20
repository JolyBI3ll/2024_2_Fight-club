package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"mime/multipart"
	"net/http"
)

type MockJwtTokenService struct {
	MockCreate            func(s *sessions.Session, tokenExpTime int64) (string, error)
	MockValidate          func(tokenString string) (*middleware.JwtCsrfClaims, error)
	MockParseSecretGetter func(token *jwt.Token) (interface{}, error)
}

func (m *MockJwtTokenService) Create(s *sessions.Session, tokenExpTime int64) (string, error) {
	return m.MockCreate(s, tokenExpTime)
}

func (m *MockJwtTokenService) Validate(tokenString string) (*middleware.JwtCsrfClaims, error) {
	return m.MockValidate(tokenString)
}

func (m *MockJwtTokenService) ParseSecretGetter(token *jwt.Token) (interface{}, error) {
	return m.MockParseSecretGetter(token)
}

type MockServiceSession struct {
	MockGetUserID      func(ctx context.Context, r *http.Request) (string, error)
	MockLogoutSession  func(ctx context.Context, r *http.Request, w http.ResponseWriter) error
	MockCreateSession  func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error)
	MockGetSessionData func(ctx context.Context, r *http.Request) (*map[string]interface{}, error)
	MockGetSession     func(ctx context.Context, r *http.Request) (*sessions.Session, error)
}

func (m *MockServiceSession) GetUserID(ctx context.Context, r *http.Request) (string, error) {
	return m.MockGetUserID(ctx, r)
}

func (m *MockServiceSession) LogoutSession(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	return m.MockLogoutSession(ctx, r, w)
}

func (m *MockServiceSession) CreateSession(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (*sessions.Session, error) {
	return m.MockCreateSession(ctx, r, w, user)
}

func (m *MockServiceSession) GetSessionData(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
	return m.MockGetSessionData(ctx, r)
}

func (m *MockServiceSession) GetSession(ctx context.Context, r *http.Request) (*sessions.Session, error) {
	return m.MockGetSession(ctx, r)
}

type MockAdUseCase struct {
	MockGetAllPlaces     func(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error)
	MockGetOnePlace      func(ctx context.Context, adId string) (domain.GetAllAdsResponse, error)
	MockCreatePlace      func(ctx context.Context, place *domain.Ad, fileHeader []*multipart.FileHeader, newPlace domain.CreateAdRequest) error
	MockUpdatePlace      func(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader []*multipart.FileHeader, updatedPlace domain.UpdateAdRequest) error
	MockDeletePlace      func(ctx context.Context, adId string, userId string) error
	MockGetPlacesPerCity func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error)
	MockGetUserPlaces    func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	MockDeleteAdImage    func(ctx context.Context, adId string, imageId int, userId string) error
}

func (m *MockAdUseCase) DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) error {
	return m.MockDeleteAdImage(ctx, adId, imageId, userId)
}

func (m *MockAdUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetAllPlaces(ctx, filter)
}

func (m *MockAdUseCase) GetOnePlace(ctx context.Context, adId string) (domain.GetAllAdsResponse, error) {
	return m.MockGetOnePlace(ctx, adId)
}

func (m *MockAdUseCase) CreatePlace(ctx context.Context, place *domain.Ad, fileHeader []*multipart.FileHeader, newPlace domain.CreateAdRequest) error {
	return m.MockCreatePlace(ctx, place, fileHeader, newPlace)
}

func (m *MockAdUseCase) UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader []*multipart.FileHeader, updatedPlace domain.UpdateAdRequest) error {
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
	MockGetAllPlaces     func(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error)
	MockGetPlaceById     func(ctx context.Context, adId string) (domain.GetAllAdsResponse, error)
	MockCreatePlace      func(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest) error
	MockSavePlace        func(ctx context.Context, ad *domain.Ad) error
	MockUpdatePlace      func(ctx context.Context, ad *domain.Ad, adId string, userId string, updatedAd domain.UpdateAdRequest) error
	MockDeletePlace      func(ctx context.Context, adId string, userId string) error
	MockGetPlacesPerCity func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error)
	MockSaveImages       func(ctx context.Context, adUUID string, imagePaths []string) error
	MockGetAdImages      func(ctx context.Context, adId string) ([]string, error)
	MockGetUserPlaces    func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	MockDeleteAdImage    func(ctx context.Context, adId string, imageId int, userId string) (string, error)
}

func (m *MockAdRepository) DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) (string, error) {
	return m.MockDeleteAdImage(ctx, adId, imageId, userId)
}

func (m *MockAdRepository) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	return m.MockGetAllPlaces(ctx, filter)
}

func (m *MockAdRepository) GetPlaceById(ctx context.Context, adId string) (domain.GetAllAdsResponse, error) {
	return m.MockGetPlaceById(ctx, adId)
}

func (m *MockAdRepository) CreatePlace(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest) error {
	return m.MockCreatePlace(ctx, ad, newAd)
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

type MockMinioService struct {
	UploadFileFunc func(file *multipart.FileHeader, path string) (string, error)
	DeleteFileFunc func(path string) error
}

func (m *MockMinioService) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	return m.UploadFileFunc(file, path)
}

func (m *MockMinioService) DeleteFile(path string) error {
	return m.DeleteFileFunc(path)
}
