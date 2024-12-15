package utils

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	"errors"
	"fmt"
	"math"
)

func ConvertAuthResponseProtoToGo(response *gen.UserResponse, userSession string) (domain.AuthResponse, error) {
	if response == nil || response.User == nil {
		return domain.AuthResponse{}, errors.New("invalid response or user nil")
	}

	return domain.AuthResponse{
		SessionId: userSession,
		User: domain.AuthData{
			Id:       response.User.Id,
			Username: response.User.Username,
			Email:    response.User.Email,
		},
	}, nil
}

func ConvertUserResponseProtoToGo(user *gen.MetadataOneUser) (domain.UserDataResponse, error) {
	if user == nil {
		return domain.UserDataResponse{}, errors.New("user response is nil")
	}

	return domain.UserDataResponse{
		Uuid:       user.Uuid,
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		Score:      math.Round(float64(user.Score)*10) / 10,
		Avatar:     user.Avatar,
		Sex:        user.Sex,
		GuestCount: int(user.GuestCount),
		Birthdate:  user.Birthdate.AsTime(),
		IsHost:     user.IsHost,
	}, nil
}

func ConvertUsersProtoToGo(users *gen.AllUsersResponse) ([]*domain.UserDataResponse, error) {
	var body []*domain.UserDataResponse

	for _, user := range users.Users {
		userResponse, err := ConvertUserResponseProtoToGo(user)
		if err != nil {
			return nil, fmt.Errorf("error converting user %s: %v", user.Uuid, err)
		}

		body = append(body, &userResponse)
	}

	return body, nil
}

func ConvertSessionDataProtoToGo(sessionData *gen.SessionDataResponse) (domain.SessionData, error) {
	if sessionData == nil {
		return domain.SessionData{}, errors.New("sessionData is nil")
	}

	return domain.SessionData{
		Id:     sessionData.Id,
		Avatar: sessionData.Avatar,
	}, nil
}
