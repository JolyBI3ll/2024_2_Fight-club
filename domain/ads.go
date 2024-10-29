package domain

import (
	"context"
	"time"
)

type Ad struct {
	UUID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:uuid" json:"id"`
	CityID          int       `gorm:"column:cityId;not null" json:"cityId"`
	AuthorUUID      string    `gorm:"column:authorUUID;not null" json:"authorUUID"`
	Address         string    `gorm:"type:varchar(255);column:address" json:"address"`
	PublicationDate time.Time `gorm:"type:date;column:publicationDate" json:"publicationDate"`
	Distance        float64   `gorm:"type:numeric;column:distance" json:"distance"`
	City            City      `gorm:"foreignKey:CityID;references:ID"`
	Author          User      `gorm:"foreignKey:AuthorUUID;references:UUID"`
}

type AdFilter struct {
	Location    string
	Rating      string
	NewThisWeek string
	HostGender  string
	GuestCount  string
}

type AdRepository interface {
	GetAllPlaces(ctx context.Context, filter AdFilter) ([]Ad, error)
	GetPlaceById(ctx context.Context, adId string) (Ad, error)
	CreatePlace(ctx context.Context, ad *Ad) error
	SavePlace(ctx context.Context, ad *Ad) error
	UpdatePlace(ctx context.Context, ad *Ad, adId string, userId string) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]Ad, error)
}
