package utils

import (
	"2024_2_FIGHT-CLUB/domain"
	adsGen "2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	authGen "2024_2_FIGHT-CLUB/microservices/auth_service/controller/gen"
	cityGen "2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

type UtilsInterface interface {
	ConvertGetAllAdsResponseProtoToGo(proto *adsGen.GetAllAdsResponseList) (domain.GetAllAdsListResponse, error)
	ConvertAdProtoToGo(ad *adsGen.GetAllAdsResponse) (domain.GetAllAdsResponse, error)
	ConvertAuthResponseProtoToGo(response *authGen.UserResponse, userSession string) (domain.AuthResponse, error)
	ConvertUserResponseProtoToGo(user *authGen.MetadataOneUser) (domain.UserDataResponse, error)
	ConvertUsersProtoToGo(users *authGen.AllUsersResponse) ([]*domain.UserDataResponse, error)
	ConvertSessionDataProtoToGo(sessionData *authGen.SessionDataResponse) (domain.SessionData, error)
	ConvertAllCitiesProtoToGo(cities *cityGen.GetCitiesResponse) ([]*domain.City, error)
	ConvertOneCityProtoToGo(city *cityGen.City) (domain.City, error)
}

type Utils struct{}

func NewUtilsInterface() UtilsInterface {
	return &Utils{}
}

const layout = "2006-01-02"

func parseDate(dateStr, adID, fieldName string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Printf("Error parsing %s for ad %s: %v\n", fieldName, adID, err)
		return time.Time{}, errors.New("error parsing date for ad")
	}
	return parsedDate, nil
}

func (u *Utils) ConvertGetAllAdsResponseProtoToGo(proto *adsGen.GetAllAdsResponseList) (domain.GetAllAdsListResponse, error) {
	ads := make([]domain.GetAllAdsResponse, 0, len(proto.Housing))

	for _, ad := range proto.Housing {
		adConverted, err := u.ConvertAdProtoToGo(ad)
		if err != nil {
			return domain.GetAllAdsListResponse{}, err
		}
		ads = append(ads, adConverted)
	}

	return domain.GetAllAdsListResponse{Housing: ads}, nil
}

func (u *Utils) ConvertAdProtoToGo(ad *adsGen.GetAllAdsResponse) (domain.GetAllAdsResponse, error) {
	if ad == nil {
		return domain.GetAllAdsResponse{}, errors.New("ad is nil")
	}

	if ad.AdAuthor == nil {
		return domain.GetAllAdsResponse{}, errors.New("adAuthor is nil")
	}

	parsedPublicationDate, err := parseDate(ad.PublicationDate, ad.Id, "PublicationDate")
	if err != nil {
		return domain.GetAllAdsResponse{}, err
	}

	parsedEndBoostDate, err := parseDate(ad.EndBoostDate, ad.Id, "EndBoostDate")
	if err != nil {
		return domain.GetAllAdsResponse{}, err
	}

	parsedDateTo, err := parseDate(ad.AdDateTo, ad.Id, "AdDateTo")
	if err != nil {
		return domain.GetAllAdsResponse{}, err
	}

	parsedDateFrom, err := parseDate(ad.AdDateFrom, ad.Id, "AdDateFrom")
	if err != nil {
		return domain.GetAllAdsResponse{}, err
	}

	parsedBirthDate, err := parseDate(ad.AdAuthor.BirthDate, ad.Id, "BirthDate")
	if err != nil {
		return domain.GetAllAdsResponse{}, err
	}

	// Преобразуем объявление
	return domain.GetAllAdsResponse{
		UUID:            ad.Id,
		CityID:          int(ad.CityId),
		AuthorUUID:      ad.AuthorUUID,
		Address:         ad.Address,
		PublicationDate: parsedPublicationDate,
		Description:     ad.Description,
		RoomsNumber:     int(ad.RoomsNumber),
		ViewsCount:      int(ad.ViewsCount),
		SquareMeters:    int(ad.SquareMeters),
		Floor:           int(ad.Floor),
		BuildingType:    ad.BuildingType,
		HasBalcony:      ad.HasBalcony,
		HasElevator:     ad.HasElevator,
		HasGas:          ad.HasGas,
		LikesCount:      int(ad.LikesCount),
		Priority:        int(ad.Priority),
		EndBoostDate:    parsedEndBoostDate,
		CityName:        ad.CityName,
		AdDateFrom:      parsedDateFrom,
		AdDateTo:        parsedDateTo,
		IsFavorite:      ad.IsFavorite,
		AdAuthor: domain.UserResponce{
			Rating:     math.Round(float64(ad.AdAuthor.Rating)*10) / 10,
			Avatar:     ad.AdAuthor.Avatar,
			Name:       ad.AdAuthor.Name,
			Sex:        ad.AdAuthor.Sex,
			Birthdate:  parsedBirthDate,
			GuestCount: int(ad.AdAuthor.GuestCount),
		},
		Images: u.convertImagesResponseProtoToGo(ad.Images),
		Rooms:  u.convertAdRoomsResponseProtoToGo(ad.Rooms),
	}, nil
}

// Вспомогательные функции для конвертации массивов
func (u *Utils) convertImagesResponseProtoToGo(protoImages []*adsGen.ImageResponse) []domain.ImageResponse {
	images := make([]domain.ImageResponse, len(protoImages))
	for i, img := range protoImages {
		images[i] = domain.ImageResponse{
			ID:        int(img.Id),
			ImagePath: img.Path,
		}
	}
	return images
}

func (u *Utils) convertAdRoomsResponseProtoToGo(protoRooms []*adsGen.AdRooms) []domain.AdRoomsResponse {
	rooms := make([]domain.AdRoomsResponse, len(protoRooms))
	for i, room := range protoRooms {
		rooms[i] = domain.AdRoomsResponse{
			Type:         room.Type,
			SquareMeters: int(room.SquareMeters),
		}
	}
	return rooms
}

func (u *Utils) ConvertAuthResponseProtoToGo(response *authGen.UserResponse, userSession string) (domain.AuthResponse, error) {
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

func (u *Utils) ConvertUserResponseProtoToGo(user *authGen.MetadataOneUser) (domain.UserDataResponse, error) {
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

func (u *Utils) ConvertUsersProtoToGo(users *authGen.AllUsersResponse) ([]*domain.UserDataResponse, error) {
	var body []*domain.UserDataResponse

	for _, user := range users.Users {
		userResponse, err := u.ConvertUserResponseProtoToGo(user)
		if err != nil {
			return nil, fmt.Errorf("error converting user %s: %v", user.Uuid, err)
		}

		body = append(body, &userResponse)
	}

	return body, nil
}

func (u *Utils) ConvertSessionDataProtoToGo(sessionData *authGen.SessionDataResponse) (domain.SessionData, error) {
	if sessionData == nil {
		return domain.SessionData{}, errors.New("sessionData is nil")
	}

	return domain.SessionData{
		Id:     sessionData.Id,
		Avatar: sessionData.Avatar,
	}, nil
}

func (u *Utils) ConvertAllCitiesProtoToGo(cities *cityGen.GetCitiesResponse) ([]*domain.City, error) {
	if cities == nil || cities.Cities == nil {
		return []*domain.City{}, errors.New("cities is nil")
	}
	var body []*domain.City
	for _, city := range cities.Cities {
		cityResponse, err := u.ConvertOneCityProtoToGo(city)
		if err != nil {
			return nil, errors.New("error convert city")
		}

		body = append(body, &cityResponse)
	}

	return body, nil
}

func (u *Utils) ConvertOneCityProtoToGo(city *cityGen.City) (domain.City, error) {
	if city == nil {
		return domain.City{}, errors.New("city response is nil")
	}

	return domain.City{
		ID:          int(city.Id),
		Title:       city.Title,
		EnTitle:     city.Entitle,
		Description: city.Description,
		Image:       city.Image,
	}, nil
}
