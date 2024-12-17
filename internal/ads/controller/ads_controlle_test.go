package controller

//import (
//	"2024_2_FIGHT-CLUB/domain"
//	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
//	"2024_2_FIGHT-CLUB/microservices/ads_service/mocks"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestAdHandler_GetAllPlaces_Success(t *testing.T) {
//	// Мокаем зависимые сервисы
//	mockGrpcClient := &mocks.MockGrpcClient{}
//	mockSessionService := &mocks.MockServiceSession{}
//
//	// Создаем экземпляр AdHandler
//	handler := &AdHandler{
//		client:         mockGrpcClient,
//		sessionService: mockSessionService,
//	}
//
//	// Создаем фейковый запрос
//	req, err := http.NewRequest("GET", "/ads?location=NewYork&rating=5", nil)
//	if err != nil {
//		t.Fatalf("could not create request: %v", err)
//	}
//
//	// Мокаем поведение для GetSessionId
//	mockSessionService.On("GetSessionId", req).Return("mockSessionID", nil)
//
//	// Мокаем запрос к gRPC API
//	expectedResponse := &gen.GetAllAdsResponseList{
//		Housing: []*gen.GetAllAdsResponse{
//			{
//				Id: "1234", CityId: 1, AuthorUUID: "user123", Address: "Test Address", PublicationDate: "2024-12-17",
//			},
//		},
//	}
//
//	mockGrpcClient.On("GetAllPlaces", mock.Anything, mock.AnythingOfType("*gen.AdFilterRequest")).Return(expectedResponse, nil)
//
//	// Мокаем функцию конвертации Proto в Go
//	mockConvertService.On("ConvertGetAllAdsResponseProtoToGo", expectedResponse).Return([]domain.GetAllAdsResponse{
//		{UUID: "1234", CityID: 1, AuthorUUID: "user123", Address: "Test Address", PublicationDate: "2024-12-17"},
//	}, nil)
//
//	// Создаем тестовый HTTP-респондер
//	rr := httptest.NewRecorder()
//
//	// Вызываем метод
//	handler.GetAllPlaces(rr, req)
//
//	// Проверяем статус код ответа
//	assert.Equal(t, http.StatusOK, rr.Code)
//
//	// Проверяем, что тело ответа соответствует ожидаемому
//	expectedBody := `[{"UUID":"1234","CityID":1,"AuthorUUID":"user123","Address":"Test Address","PublicationDate":"2024-12-17"}]`
//	assert.JSONEq(t, expectedBody, rr.Body.String())
//
//	// Проверяем, что все ожидаемые вызовы были выполнены
//	mockGrpcClient.AssertExpectations(t)
//	mockSessionService.AssertExpectations(t)
//	mockConvertService.AssertExpectations(t)
//}
