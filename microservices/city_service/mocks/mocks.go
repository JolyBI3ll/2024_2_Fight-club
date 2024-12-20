package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockCitiesRepository struct {
	MockGetCities       func(ctx context.Context) ([]domain.City, error)
	MockGetCityByEnName func(ctx context.Context, cityEnName string) (domain.City, error)
}

func (m *MockCitiesRepository) GetCities(ctx context.Context) ([]domain.City, error) {
	return m.MockGetCities(ctx)
}

func (m *MockCitiesRepository) GetCityByEnName(ctx context.Context, cityEnName string) (domain.City, error) {
	return m.MockGetCityByEnName(ctx, cityEnName)
}

type MockCitiesUseCase struct {
	MockGetCities  func(ctx context.Context) ([]domain.City, error)
	MockGetOneCity func(ctx context.Context, cityEnName string) (domain.City, error)
}

func (m *MockCitiesUseCase) GetCities(ctx context.Context) ([]domain.City, error) {
	return m.MockGetCities(ctx)
}

func (m *MockCitiesUseCase) GetOneCity(ctx context.Context, cityEnName string) (domain.City, error) {
	return m.MockGetOneCity(ctx, cityEnName)
}

type MockGrpcClient struct {
	mock.Mock
}

func (m *MockGrpcClient) GetCities(ctx context.Context, in *gen.GetCitiesRequest, opts ...grpc.CallOption) (*gen.GetCitiesResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.GetCitiesResponse), args.Error(1)
}
func (m *MockGrpcClient) GetOneCity(ctx context.Context, in *gen.GetOneCityRequest, opts ...grpc.CallOption) (*gen.GetOneCityResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*gen.GetOneCityResponse), args.Error(1)
}
