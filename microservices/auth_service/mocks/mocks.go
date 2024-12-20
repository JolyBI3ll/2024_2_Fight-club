package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
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
	MockGetSessionData func(ctx context.Context, sessionID string) (*domain.SessionData, error)
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

func (m *MockServiceSession) GetSessionData(ctx context.Context, sessionID string) (*domain.SessionData, error) {
	return m.MockGetSessionData(ctx, sessionID)
}

type MockAuthUseCase struct {
	MockRegisterUser      func(ctx context.Context, creds *domain.User) error
	MockLoginUser         func(ctx context.Context, creds *domain.User) (*domain.User, error)
	MockPutUser           func(ctx context.Context, creds *domain.User, userID string, avatar []byte) error
	MockGetAllUser        func(ctx context.Context) ([]domain.User, error)
	MockGetUserById       func(ctx context.Context, userID string) (*domain.User, error)
	MockUpdateUserRegions func(ctx context.Context, regions domain.UpdateUserRegion, userId string) error
	MockDeleteUserRegion  func(ctx context.Context, regionName string, userID string) error
}

func (m *MockAuthUseCase) RegisterUser(ctx context.Context, creds *domain.User) error {
	return m.MockRegisterUser(ctx, creds)
}

func (m *MockAuthUseCase) LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error) {
	return m.MockLoginUser(ctx, creds)
}

func (m *MockAuthUseCase) PutUser(ctx context.Context, creds *domain.User, userID string, avatar []byte) error {
	return m.MockPutUser(ctx, creds, userID, avatar)
}

func (m *MockAuthUseCase) GetAllUser(ctx context.Context) ([]domain.User, error) {
	return m.MockGetAllUser(ctx)
}

func (m *MockAuthUseCase) UpdateUserRegions(ctx context.Context, regions domain.UpdateUserRegion, userId string) error {
	return m.MockUpdateUserRegions(ctx, regions, userId)
}

func (m *MockAuthUseCase) DeleteUserRegion(ctx context.Context, regionName string, userID string) error {
	return m.MockDeleteUserRegion(ctx, regionName, userID)
}

func (m *MockAuthUseCase) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	return m.MockGetUserById(ctx, userID)
}

type MockAuthRepository struct {
	GetUserByNameFunc    func(ctx context.Context, username string) (*domain.User, error)
	CreateUserFunc       func(ctx context.Context, user *domain.User) error
	SaveUserFunc         func(ctx context.Context, user *domain.User) error
	PutUserFunc          func(ctx context.Context, user *domain.User, userID string) error
	GetAllUserFunc       func(ctx context.Context) ([]domain.User, error)
	GetUserByIdFunc      func(ctx context.Context, userID string) (*domain.User, error)
	MockGetUserByEmail   func(ctx context.Context, email string) (*domain.User, error)
	MockUpdateUserRegion func(ctx context.Context, region domain.UpdateUserRegion, userId string) error
	MockDeleteUserRegion func(ctx context.Context, regionName string, userId string) error
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.MockGetUserByEmail(ctx, email)
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

func (m *MockAuthRepository) UpdateUserRegion(ctx context.Context, region domain.UpdateUserRegion, userId string) error {
	return m.MockUpdateUserRegion(ctx, region, userId)
}

func (m *MockAuthRepository) DeleteUserRegion(ctx context.Context, regionName string, userId string) error {
	return m.MockDeleteUserRegion(ctx, regionName, userId)
}

type MockMinioService struct {
	UploadFileFunc func(file []byte, contentType string, id string) (string, error)
	DeleteFileFunc func(path string) error
}

func (m *MockMinioService) UploadFile(file []byte, contentType string, id string) (string, error) {
	return m.UploadFileFunc(file, contentType, id)
}

func (m *MockMinioService) DeleteFile(path string) error {
	return m.DeleteFileFunc(path)
}

type MockGrpcClient struct {
	mock.Mock
}

func (m *MockGrpcClient) RegisterUser(ctx context.Context, in *gen.RegisterUserRequest, opts ...grpc.CallOption) (*gen.UserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.UserResponse), args.Error(1)
}
func (m *MockGrpcClient) LoginUser(ctx context.Context, in *gen.LoginUserRequest, opts ...grpc.CallOption) (*gen.UserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.UserResponse), args.Error(1)
}
func (m *MockGrpcClient) LogoutUser(ctx context.Context, in *gen.LogoutRequest, opts ...grpc.CallOption) (*gen.LogoutUserResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.LogoutUserResponse), args.Error(1)
}
func (m *MockGrpcClient) PutUser(ctx context.Context, in *gen.PutUserRequest, opts ...grpc.CallOption) (*gen.UpdateResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.UpdateResponse), args.Error(1)
}
func (m *MockGrpcClient) GetUserById(ctx context.Context, in *gen.GetUserByIdRequest, opts ...grpc.CallOption) (*gen.GetUserByIdResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.GetUserByIdResponse), args.Error(1)
}
func (m *MockGrpcClient) GetAllUsers(ctx context.Context, in *gen.Empty, opts ...grpc.CallOption) (*gen.AllUsersResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.AllUsersResponse), args.Error(1)
}
func (m *MockGrpcClient) GetSessionData(ctx context.Context, in *gen.GetSessionDataRequest, opts ...grpc.CallOption) (*gen.SessionDataResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.SessionDataResponse), args.Error(1)
}
func (m *MockGrpcClient) RefreshCsrfToken(ctx context.Context, in *gen.RefreshCsrfTokenRequest, opts ...grpc.CallOption) (*gen.RefreshCsrfTokenResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.RefreshCsrfTokenResponse), args.Error(1)
}

func (m *MockGrpcClient) UpdateUserRegions(ctx context.Context, in *gen.UpdateUserRegionsRequest, opts ...grpc.CallOption) (*gen.UpdateResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.UpdateResponse), args.Error(1)
}

func (m *MockGrpcClient) DeleteUserRegions(ctx context.Context, in *gen.DeleteUserRegionsRequest, opts ...grpc.CallOption) (*gen.UpdateResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.UpdateResponse), args.Error(1)
}
