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

type MockAuthUseCase struct {
	MockRegisterUser func(ctx context.Context, creds *domain.User, avatar *multipart.FileHeader) error
	MockLoginUser    func(ctx context.Context, creds *domain.User) (*domain.User, error)
	MockPutUser      func(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error
	MockGetAllUser   func(ctx context.Context) ([]domain.User, error)
	MockGetUserById  func(ctx context.Context, userID string) (*domain.User, error)
}

func (m *MockAuthUseCase) RegisterUser(ctx context.Context, creds *domain.User, avatar *multipart.FileHeader) error {
	return m.MockRegisterUser(ctx, creds, avatar)
}

func (m *MockAuthUseCase) LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error) {
	return m.MockLoginUser(ctx, creds)
}

func (m *MockAuthUseCase) PutUser(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error {
	return m.MockPutUser(ctx, creds, userID, avatar)
}

func (m *MockAuthUseCase) GetAllUser(ctx context.Context) ([]domain.User, error) {
	return m.MockGetAllUser(ctx)
}

func (m *MockAuthUseCase) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	return m.MockGetUserById(ctx, userID)
}

type MockAuthRepository struct {
	GetUserByNameFunc func(ctx context.Context, username string) (*domain.User, error)
	CreateUserFunc    func(ctx context.Context, user *domain.User) error
	SaveUserFunc      func(ctx context.Context, user *domain.User) error
	PutUserFunc       func(ctx context.Context, user *domain.User, userID string) error
	GetAllUserFunc    func(ctx context.Context) ([]domain.User, error)
	GetUserByIdFunc   func(ctx context.Context, userID string) (*domain.User, error)
}

func (m *MockAuthRepository) GetUserByName(ctx context.Context, username string) (*domain.User, error) {
	return m.GetUserByNameFunc(ctx, username)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return m.CreateUserFunc(ctx, user)
}

func (m *MockAuthRepository) SaveUser(ctx context.Context, user *domain.User) error {
	return m.SaveUserFunc(ctx, user)
}

func (m *MockAuthRepository) PutUser(ctx context.Context, user *domain.User, userID string) error {
	return m.PutUserFunc(ctx, user, userID)
}

func (m *MockAuthRepository) GetAllUser(ctx context.Context) ([]domain.User, error) {
	return m.GetAllUserFunc(ctx)
}

func (m *MockAuthRepository) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	return m.GetUserByIdFunc(ctx, userID)
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
