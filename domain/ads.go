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
	Description     string    `gorm:"type:text;size:1000;column:description" json:"description"`
	RoomsNumber     int       `gorm:"column:roomsNumber" json:"roomsNumber"`
	City            City      `gorm:"foreignKey:CityID;references:ID" json:"-"`
	Author          User      `gorm:"foreignKey:AuthorUUID;references:UUID" json:"-"`
}

type GetAllAdsResponse struct {
	UUID            string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:uuid" json:"id"`
	CityID          int             `gorm:"column:cityId;not null" json:"cityId"`
	AuthorUUID      string          `gorm:"column:authorUUID;not null" json:"authorUUID"`
	Address         string          `gorm:"type:varchar(255);column:address" json:"address"`
	PublicationDate time.Time       `gorm:"type:date;column:publicationDate" json:"publicationDate"`
	Description     string          `gorm:"type:text;size:1000;column:description" json:"description"`
	RoomsNumber     int             `gorm:"column:roomsNumber" json:"roomsNumber"`
	City            City            `gorm:"foreignKey:CityID;references:ID" json:"-"`
	Author          User            `gorm:"foreignKey:AuthorUUID;references:UUID" json:"-"`
	Cityname        string          `json:"cityName"`
	AdAuthor        UserResponce    `gorm:"-" json:"author"`
	Images          []ImageResponse `gorm:"-" json:"images"`
}

type CreateAdRequest struct {
	CityName    string `form:"cityName" json:"cityName"`
	Address     string `form:"address" json:"address"`
	Description string `form:"description" json:"description"`
	RoomsNumber int    `form:"roomsNumber" json:"roomsNumber"`
}

type UpdateAdRequest struct {
	CityName    string `form:"cityName" json:"cityName"`
	Address     string `form:"address" json:"address"`
	Description string `form:"description" json:"description"`
	RoomsNumber int    `form:"roomsNumber" json:"roomsNumber"`
}

type AdFilter struct {
	Location    string
	Rating      string
	NewThisWeek string
	HostGender  string
	GuestCount  string
}

type AdRepository interface {
	GetAllPlaces(ctx context.Context, filter AdFilter) ([]GetAllAdsResponse, error)
	GetPlaceById(ctx context.Context, adId string) (GetAllAdsResponse, error)
	CreatePlace(ctx context.Context, ad *Ad, newAd CreateAdRequest) error
	SavePlace(ctx context.Context, ad *Ad) error
	UpdatePlace(ctx context.Context, ad *Ad, adId string, userId string, updatedAd UpdateAdRequest) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]GetAllAdsResponse, error)
	SaveImages(ctx context.Context, adUUID string, imagePaths []string) error
	GetAdImages(ctx context.Context, adId string) ([]string, error)
	GetUserPlaces(ctx context.Context, userId string) ([]GetAllAdsResponse, error)
	DeleteAdImage(ctx context.Context, adId string, imageId string, userId string) (string, error)
}
