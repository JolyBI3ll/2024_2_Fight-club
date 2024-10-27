package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
	"mime/multipart"
	"net/http"
)

type MockServiceSession struct {
	MockGetUserID      func(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error)
	MockLogoutSession  func(ctx context.Context, r *http.Request, w http.ResponseWriter) error
	MockCreateSession  func(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (string, error)
	MockGetSessionData func(ctx context.Context, r *http.Request) (*map[string]interface{}, error)
}

func (m *MockServiceSession) GetUserID(ctx context.Context, r *http.Request, w http.ResponseWriter) (string, error) {
	return m.MockGetUserID(ctx, r, w)
}

func (m *MockServiceSession) LogoutSession(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	return m.MockLogoutSession(ctx, r, w)
}

func (m *MockServiceSession) CreateSession(ctx context.Context, r *http.Request, w http.ResponseWriter, user *domain.User) (string, error) {
	return m.MockCreateSession(ctx, r, w, user)
}

func (m *MockServiceSession) GetSessionData(ctx context.Context, r *http.Request) (*map[string]interface{}, error) {
	return m.MockGetSessionData(ctx, r)
}

type MockAdUseCase struct {
	MockGetAllPlaces     func(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error)
	MockGetOnePlace      func(ctx context.Context, adId string) (domain.Ad, error)
	MockCreatePlace      func(ctx context.Context, ad *domain.Ad, files []*multipart.FileHeader) error
	MockUpdatePlace      func(ctx context.Context, ad *domain.Ad, adId string, userId string, files []*multipart.FileHeader) error
	MockDeletePlace      func(ctx context.Context, adId string, userId string) error
	MockGetPlacesPerCity func(ctx context.Context, city string) ([]domain.Ad, error)
}

func (m *MockAdUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
	return m.MockGetAllPlaces(ctx, filter)
}

func (m *MockAdUseCase) GetOnePlace(ctx context.Context, adId string) (domain.Ad, error) {
	return m.MockGetOnePlace(ctx, adId)
}

func (m *MockAdUseCase) CreatePlace(ctx context.Context, ad *domain.Ad, files []*multipart.FileHeader) error {
	return m.MockCreatePlace(ctx, ad, files)
}

func (m *MockAdUseCase) UpdatePlace(ctx context.Context, ad *domain.Ad, adId string, userId string, files []*multipart.FileHeader) error {
	return m.MockUpdatePlace(ctx, ad, adId, userId, files)
}

func (m *MockAdUseCase) DeletePlace(ctx context.Context, adId string, userId string) error {
	return m.MockDeletePlace(ctx, adId, userId)
}

func (m *MockAdUseCase) GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error) {
	return m.MockGetPlacesPerCity(ctx, city)
}

type MockAdRepository struct {
	GetAllPlacesFunc     func(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error)
	GetPlaceByIdFunc     func(ctx context.Context, adId string) (domain.Ad, error)
	CreatePlaceFunc      func(ctx context.Context, place *domain.Ad) error
	SavePlaceFunc        func(ctx context.Context, place *domain.Ad) error
	UpdatePlaceFunc      func(ctx context.Context, place *domain.Ad, adId, userId string) error
	DeletePlaceFunc      func(ctx context.Context, adId, userId string) error
	GetPlacesPerCityFunc func(ctx context.Context, city string) ([]domain.Ad, error)
}

func (m *MockAdRepository) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
	return m.GetAllPlacesFunc(ctx, filter)
}

func (m *MockAdRepository) GetPlaceById(ctx context.Context, adId string) (domain.Ad, error) {
	return m.GetPlaceByIdFunc(ctx, adId)
}

func (m *MockAdRepository) CreatePlace(ctx context.Context, place *domain.Ad) error {
	return m.CreatePlaceFunc(ctx, place)
}

func (m *MockAdRepository) SavePlace(ctx context.Context, place *domain.Ad) error {
	return m.SavePlaceFunc(ctx, place)
}

func (m *MockAdRepository) UpdatePlace(ctx context.Context, place *domain.Ad, adId, userId string) error {
	return m.UpdatePlaceFunc(ctx, place, adId, userId)
}

func (m *MockAdRepository) DeletePlace(ctx context.Context, adId, userId string) error {
	return m.DeletePlaceFunc(ctx, adId, userId)
}

func (m *MockAdRepository) GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error) {
	return m.GetPlacesPerCityFunc(ctx, city)
}

// MockMinioService - структура для ручного мока minioServiceInterface.
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
