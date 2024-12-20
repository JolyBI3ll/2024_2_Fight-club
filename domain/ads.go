package domain

//go:generate easyjson -all ads.go

import (
	"context"
	"time"
)

//easyjson:json
type PlacesResponse struct {
	Places GetAllAdsListResponse `json:"places"`
}

//easyjson:json
type GetAllAdsListResponse struct {
	Housing []GetAllAdsResponse `json:"housing"`
}

//easyjson:json
type GetOneAdResponse struct {
	Place GetAllAdsResponse `json:"place"`
}

//easyjson:json
type Ad struct {
	UUID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:uuid" json:"id"`
	CityID          int       `gorm:"column:cityId;not null" json:"cityId"`
	AuthorUUID      string    `gorm:"column:authorUUID;not null" json:"authorUUID"`
	Address         string    `gorm:"type:varchar(255);column:address" json:"address"`
	PublicationDate time.Time `gorm:"type:date;column:publicationDate" json:"publicationDate"`
	Description     string    `gorm:"type:text;size:1000;column:description" json:"description"`
	RoomsNumber     int       `gorm:"column:roomsNumber" json:"roomsNumber"`
	ViewsCount      int       `gorm:"column:viewsCount;default:0" json:"viewsCount"`
	SquareMeters    int       `gorm:"column:squareMeters" json:"squareMeters"`
	Floor           int       `gorm:"column:floor" json:"floor"`
	BuildingType    string    `gorm:"type:text;size:255;column:buildingType" json:"buildingType"`
	HasBalcony      bool      `gorm:"type:bool;default:false;column:hasBalcony" json:"hasBalcony"`
	HasElevator     bool      `gorm:"type:bool;default:false;column:hasElevator" json:"hasElevator"`
	HasGas          bool      `gorm:"type:bool;default:false;column:hasGas" json:"hasGas"`
	LikesCount      int       `gorm:"column:likesCount;default:0" json:"likesCount"`
	Priority        int       `gorm:"column:priority;default:0" json:"priority"`
	EndBoostDate    time.Time `gorm:"type:date;column:endBoostDate" json:"endBoostDate"`
	City            City      `gorm:"foreignKey:CityID;references:ID" json:"-"`
	Author          User      `gorm:"foreignKey:AuthorUUID;references:UUID" json:"-"`
}

type Favorites struct {
	AdId   string `gorm:"primaryKey;column:adId" json:"adId"`
	UserId string `gorm:"primaryKey;column:userId" json:"userId"`
	User   User   `gorm:"foreignKey:UserId;references:UUID" json:"-"`
	Ad     Ad     `gorm:"foreignKey:AdId;references:UUID" json:"-"`
}

type GetAllAdsResponse struct {
	UUID            string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:uuid" json:"id"`
	CityID          int               `gorm:"column:cityId;not null" json:"cityId"`
	AuthorUUID      string            `gorm:"column:authorUUID;not null" json:"authorUUID"`
	Address         string            `gorm:"type:varchar(255);column:address" json:"address"`
	PublicationDate time.Time         `gorm:"type:date;column:publicationDate" json:"publicationDate"`
	Description     string            `gorm:"type:text;size:1000;column:description" json:"description"`
	RoomsNumber     int               `gorm:"column:roomsNumber" json:"roomsNumber"`
	City            City              `gorm:"foreignKey:CityID;references:ID" json:"-"`
	Author          User              `gorm:"foreignKey:AuthorUUID;references:UUID" json:"-"`
	ViewsCount      int               `gorm:"column:viewsCount;default:0" json:"viewsCount"`
	SquareMeters    int               `gorm:"column:squareMeters" json:"squareMeters"`
	Floor           int               `gorm:"column:floor" json:"floor"`
	BuildingType    string            `gorm:"type:text;size:255;column:buildingType" json:"buildingType"`
	HasBalcony      bool              `gorm:"type:bool;default:false;column:hasBalcony" json:"hasBalcony"`
	HasElevator     bool              `gorm:"type:bool;default:false;column:hasElevator" json:"hasElevator"`
	HasGas          bool              `gorm:"type:bool;default:false;column:hasGas" json:"hasGas"`
	LikesCount      int               `gorm:"column:likesCount;default:0" json:"likesCount"`
	Priority        int               `gorm:"column:priority;default:0" json:"priority"`
	EndBoostDate    time.Time         `gorm:"type:date;column:endBoostDate" json:"endBoostDate"`
	CityName        string            `json:"cityName"`
	AdDateFrom      time.Time         `json:"adDateFrom"`
	AdDateTo        time.Time         `json:"adDateTo"`
	IsFavorite      bool              `json:"isFavorite"`
	AdAuthor        UserResponce      `gorm:"-" json:"author"`
	Images          []ImageResponse   `gorm:"-" json:"images"`
	Rooms           []AdRoomsResponse `gorm:"-" json:"rooms"`
}

type CreateAdRequest struct {
	CityName     string            `form:"cityName" json:"cityName"`
	Address      string            `form:"address" json:"address"`
	Description  string            `form:"description" json:"description"`
	RoomsNumber  int               `form:"roomsNumber" json:"roomsNumber"`
	DateFrom     time.Time         `form:"dateFrom" json:"dateFrom"`
	DateTo       time.Time         `form:"dateTo" json:"dateTo"`
	Rooms        []AdRoomsResponse `form:"rooms" json:"rooms"`
	SquareMeters int               `form:"squareMeters" json:"squareMeters"`
	Floor        int               `form:"floor" json:"floor"`
	BuildingType string            `form:"buildingType" json:"buildingType"`
	HasBalcony   bool              `form:"hasBalcony" json:"hasBalcony"`
	HasElevator  bool              `form:"hasElevator" json:"hasElevator"`
	HasGas       bool              `form:"hasGas" json:"hasGas"`
}

type UpdateAdRequest struct {
	CityName     string            `form:"cityName" json:"cityName"`
	Address      string            `form:"address" json:"address"`
	Description  string            `form:"description" json:"description"`
	RoomsNumber  int               `form:"roomsNumber" json:"roomsNumber"`
	DateFrom     time.Time         `form:"dateFrom" json:"dateFrom"`
	DateTo       time.Time         `form:"dateTo" json:"dateTo"`
	Rooms        []AdRoomsResponse `form:"rooms" json:"rooms"`
	SquareMeters int               `form:"squareMeters" json:"squareMeters"`
	Floor        int               `form:"floor" json:"floor"`
	BuildingType string            `form:"buildingType" json:"buildingType"`
	HasBalcony   bool              `form:"hasBalcony" json:"hasBalcony"`
	HasElevator  bool              `form:"hasElevator" json:"hasElevator"`
	HasGas       bool              `form:"hasGas" json:"hasGas"`
}

type AdFilter struct {
	Location    string
	Rating      string
	NewThisWeek string
	HostGender  string
	GuestCount  string
	Limit       int
	Offset      int
	DateFrom    time.Time
	DateTo      time.Time
	Favorites   string
}

type PaymentInfo struct {
	CardNumber     string `form:"cardNumber" json:"cardNumber"`
	CardExpiry     string `form:"cardExpiry" json:"cardExpiry"`
	CardCvc        string `form:"cardCVC" json:"cardCVC"`
	DonationAmount string `form:"donationAmount" json:"donationAmount"`
}

type AdRepository interface {
	GetAllPlaces(ctx context.Context, filter AdFilter, userId string) ([]GetAllAdsResponse, error)
	GetPlaceById(ctx context.Context, adId string) (GetAllAdsResponse, error)
	CreatePlace(ctx context.Context, ad *Ad, newAd CreateAdRequest, userId string) error
	UpdatePlace(ctx context.Context, ad *Ad, adId string, userId string, updatedAd UpdateAdRequest) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]GetAllAdsResponse, error)
	SaveImages(ctx context.Context, adUUID string, imagePaths []string) error
	GetAdImages(ctx context.Context, adId string) ([]string, error)
	GetUserPlaces(ctx context.Context, userId string) ([]GetAllAdsResponse, error)
	DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) (string, error)
	UpdateViewsCount(ctx context.Context, ad GetAllAdsResponse) (GetAllAdsResponse, error)
	AddToFavorites(ctx context.Context, adId string, userId string) error
	DeleteFromFavorites(ctx context.Context, adId string, userId string) error
	GetUserFavorites(ctx context.Context, userId string) ([]GetAllAdsResponse, error)
	UpdateFavoritesCount(ctx context.Context, adId string) error
	UpdatePriority(ctx context.Context, adId string, userId string, amount int) error
	ResetExpiredPriorities(ctx context.Context) error
}
