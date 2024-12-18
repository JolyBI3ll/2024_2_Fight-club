package utils

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	authGen "2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	cityGen "2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"github.com/stretchr/testify/mock"
)

type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) ConvertGetAllAdsResponseProtoToGo(proto *gen.GetAllAdsResponseList) (domain.GetAllAdsListResponse, error) {
	args := m.Called(proto)
	if res, ok := args.Get(0).(domain.GetAllAdsListResponse); ok {
		return res, args.Error(1)
	}
	return domain.GetAllAdsListResponse{}, args.Error(1)
}

func (m *MockUtils) ConvertAdProtoToGo(ad *gen.GetAllAdsResponse) (domain.GetAllAdsResponse, error) {
	args := m.Called(ad)
	if res, ok := args.Get(0).(domain.GetAllAdsResponse); ok {
		return res, args.Error(1)
	}
	return domain.GetAllAdsResponse{}, args.Error(1)
}

func (m *MockUtils) ConvertAuthResponseProtoToGo(response *authGen.UserResponse, userSession string) (domain.AuthResponse, error) {
	args := m.Called(response, userSession)
	if res, ok := args.Get(0).(domain.AuthResponse); ok {
		return res, args.Error(1)
	}
	return domain.AuthResponse{}, args.Error(1)
}

func (m *MockUtils) ConvertUserResponseProtoToGo(user *authGen.MetadataOneUser) (domain.UserDataResponse, error) {
	args := m.Called(user)
	if res, ok := args.Get(0).(domain.UserDataResponse); ok {
		return res, args.Error(1)
	}
	return domain.UserDataResponse{}, args.Error(1)
}

func (m *MockUtils) ConvertUsersProtoToGo(users *authGen.AllUsersResponse) ([]*domain.UserDataResponse, error) {
	args := m.Called(users)
	if res, ok := args.Get(0).([]*domain.UserDataResponse); ok {
		return res, args.Error(1)
	}
	return []*domain.UserDataResponse{}, args.Error(1)
}

func (m *MockUtils) ConvertSessionDataProtoToGo(sessionData *authGen.SessionDataResponse) (domain.SessionData, error) {
	args := m.Called(sessionData)
	if res, ok := args.Get(0).(domain.SessionData); ok {
		return res, args.Error(1)
	}
	return domain.SessionData{}, args.Error(1)
}

func (m *MockUtils) ConvertAllCitiesProtoToGo(cities *cityGen.GetCitiesResponse) ([]*domain.City, error) {
	args := m.Called(cities)
	var trueRes []*domain.City
	if res, ok := args.Get(0).(domain.AllCitiesResponse); ok {
		for _, city := range res.Cities {
			trueRes = append(trueRes, city)
		}
		return trueRes, args.Error(1)
	}
	return []*domain.City{}, args.Error(1)
}

func (m *MockUtils) ConvertOneCityProtoToGo(city *cityGen.City) (domain.City, error) {
	args := m.Called(city)
	if res, ok := args.Get(0).(domain.City); ok {
		return res, args.Error(1)
	}
	return domain.City{}, args.Error(1)
}
