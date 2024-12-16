package utils

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/microservices/city_service/controller/gen"
	"errors"
)

func ConvertAllCitiesProtoToGo(cities *gen.GetCitiesResponse) ([]*domain.City, error) {
	if cities == nil || cities.Cities == nil {
		return []*domain.City{}, errors.New("cities is nil")
	}
	var body []*domain.City
	for _, city := range cities.Cities {
		cityResponse, err := ConvertOneCityProtoToGo(city)
		if err != nil {
			return nil, errors.New("error convert city")
		}

		body = append(body, &cityResponse)
	}

	return body, nil
}

func ConvertOneCityProtoToGo(city *gen.City) (domain.City, error) {
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
