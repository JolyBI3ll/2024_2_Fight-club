package domain

//go:generate easyjson -all auth.go

import (
	"context"
	"time"
)

//easyjson:json
type CSRFTokenResponse struct {
	Token string `json:"csrf_token"`
}

//easyjson:json
type SessionData struct {
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
}

//easyjson:json
type GetAllUsersResponse struct {
	Users []*UserDataResponse `json:"users"`
}

//easyjson:json
type AuthResponse struct {
	SessionId string   `json:"session_id"`
	User      AuthData `json:"user"`
}

//easyjson:json
type UpdateUserRegion struct {
	RegionName       string `json:"regionName"`
	StartVisitedDate string `json:"startVisitedDate"`
	EndVisitedDate   string `json:"endVisitedDate"`
}

type AuthData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

//easyjson:json
type User struct {
	UUID       string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:uuid" json:"id"`
	Username   string    `gorm:"type:varchar(20);unique;not null;column:username" json:"username"`
	Password   string    `gorm:"type:varchar(255);not null;column:password" json:"password"`
	Email      string    `gorm:"type:varchar(255);unique;not null;column:email" json:"email"`
	Name       string    `gorm:"type:varchar(50);not null;column:name" json:"name"`
	Score      float64   `gorm:"type:numeric;column:score" json:"score"`
	Avatar     string    `gorm:"type:text;column:avatar;size:1000;default:/images/default.png" json:"avatar"`
	Sex        string    `gorm:"type:varchar(1);column:sex" json:"sex"`
	GuestCount int       `gorm:"column:guestCount" json:"guestCount"`
	Birthdate  time.Time `gorm:"type:date;column:birthDate" json:"birthDate"`
	IsHost     bool      `gorm:"type:boolean;default:false;column:isHost" form:"isHost" json:"isHost"`
}

type UserResponce struct {
	Rating     float64   `json:"rating"`
	Avatar     string    `json:"avatar"`
	Name       string    `json:"name"`
	Sex        string    `json:"sex"`
	Birthdate  time.Time `json:"birthDate"`
	GuestCount int       `json:"guestCount"`
}

//easyjson:json
type UserDataResponse struct {
	Uuid       string    `json:"uuid"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Score      float64   `json:"score"`
	Avatar     string    `json:"avatar"`
	Sex        string    `json:"sex"`
	GuestCount int       `json:"guestCount"`
	Birthdate  time.Time `json:"birthdate"`
	IsHost     bool      `json:"isHost"`
}

type AuthRepository interface {
	CreateUser(ctx context.Context, creds *User) error
	SaveUser(ctx context.Context, creds *User) error
	PutUser(ctx context.Context, creds *User, userID string) error
	GetAllUser(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, userID string) (*User, error)
	GetUserByName(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUserRegion(ctx context.Context, region UpdateUserRegion, userId string) error
	DeleteUserRegion(ctx context.Context, regionName string, userId string) error
}
