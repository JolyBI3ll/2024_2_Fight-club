package domain

import (
	"2024_2_FIGHT-CLUB/internal/service/type"
)

type Ad struct {
	ID              string             `gorm:"primaryKey" json:"id"`
	LocationMain    string             `json:"location_main"`
	LocationStreet  string             `json:"location_street"`
	Position        ntype.Float64Array `gorm:"type:float[]" json:"position"`
	Images          ntype.StringArray  `gorm:"type:text[]"`
	AuthorUUID      string             `json:"author_uuid"`
	PublicationDate string             `json:"publication_date"`
	AvailableDates  []string           `gorm:"type:text[]" json:"available_dates"`
	Distance        float32            `json:"distance"`
	Requests        []Request          `gorm:"foreignKey:AdID" json:"requests"`
}

type AdRepository interface {
	GetAllPlaces() ([]Ad, error)
	GetPlaceById(adId string) (Ad, error)
	CreatePlace(ad *Ad) error
	SavePlace(ad *Ad) error
	UpdatePlace(ad *Ad, adId string, userId string) error
	DeletePlace(adId string, userId string) error
}
