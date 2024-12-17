package mocks

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"github.com/stretchr/testify/mock"
	"time"
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

func (m *MockUtils) ParseDate(dateStr, adID, fieldName string) (time.Time, error) {
	args := m.Called(dateStr, adID, fieldName)
	if res, ok := args.Get(0).(time.Time); ok {
		return res, args.Error(1)
	}
	return time.Time{}, args.Error(1)
}
