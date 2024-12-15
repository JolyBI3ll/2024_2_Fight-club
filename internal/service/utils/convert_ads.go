package utils

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"errors"
	"fmt"
	"log"
	"time"
)

const layout = "2006-01-02"

func parseDate(dateStr, adID, fieldName string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Printf("Error parsing %s for ad %s: %v\n", fieldName, adID, err)
		return time.Time{}, errors.New(fmt.Sprintf("Error parsing %s for ad %s: %v", fieldName, adID, err))
	}
	return parsedDate, nil
}

func ConvertGetAllAdsResponseProtoToGo(proto *gen.GetAllAdsResponseList) (domain.GetAllAdsListResponse, error) {
	ads := make([]domain.GetAllAdsResponse, 0, len(proto.Housing))

	for _, ad := range proto.Housing {
		adConverted, err := ConvertAdProtoToGo(ad)
		if err != nil {
			return domain.GetAllAdsListResponse{}, err
		}
		ads = append(ads, adConverted)
	}

	return domain.GetAllAdsListResponse{Housing: ads}, nil
}

func ConvertAdProtoToGo(ad *gen.GetAllAdsResponse) (domain.GetAllAdsResponse, error) {
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
			Rating:     float64(ad.AdAuthor.Rating),
			Avatar:     ad.AdAuthor.Avatar,
			Name:       ad.AdAuthor.Name,
			Sex:        ad.AdAuthor.Sex,
			Birthdate:  parsedBirthDate,
			GuestCount: int(ad.AdAuthor.GuestCount),
		},
		Images: convertImagesResponseProtoToGo(ad.Images),
		Rooms:  convertAdRoomsResponseProtoToGo(ad.Rooms),
	}, nil
}

// Вспомогательные функции для конвертации массивов
func convertImagesResponseProtoToGo(protoImages []*gen.ImageResponse) []domain.ImageResponse {
	images := make([]domain.ImageResponse, len(protoImages))
	for i, img := range protoImages {
		images[i] = domain.ImageResponse{
			ID:        int(img.Id),
			ImagePath: img.Path,
		}
	}
	return images
}

func convertAdRoomsResponseProtoToGo(protoRooms []*gen.AdRooms) []domain.AdRoomsResponse {
	rooms := make([]domain.AdRoomsResponse, len(protoRooms))
	for i, room := range protoRooms {
		rooms[i] = domain.AdRoomsResponse{
			Type:         room.Type,
			SquareMeters: int(room.SquareMeters),
		}
	}
	return rooms
}
