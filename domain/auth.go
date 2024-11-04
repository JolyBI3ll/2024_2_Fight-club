package domain

import (
	"context"
	"time"
)

type User struct {
	UUID       string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:uuid" json:"id"`
	Username   string    `gorm:"type:varchar(20);unique;not null;column:username" json:"username"`
	Password   string    `gorm:"type:varchar(255);not null;column:password" json:"password"`
	Email      string    `gorm:"type:varchar(255);unique;not null;column:email" json:"email"`
	Name       string    `gorm:"type:varchar(50);not null;column:name" json:"name"`
	Score      float64   `gorm:"type:numeric;column:score" json:"score"`
	Avatar     string    `gorm:"type:text;column:avatar;size:1000;default:images/default.png" json:"avatar"`
	Sex        string    `gorm:"type:varchar(1);column:sex" json:"sex"`
	GuestCount int       `gorm:"column:guestCount" json:"guestCount"`
	Birthdate  time.Time `gorm:"type:date;column:birthDate" json:"birthDate"`
	IsHost     bool      `gorm:"type:boolean;default:false;column:isHost" form:"isHost" json:"isHost"`
}

type UserResponce struct {
	Rating float64 `json:"rating"`
	Avatar string  `json:"avatar"`
	Name   string  `json:"name"`
}

type AuthRepository interface {
	CreateUser(ctx context.Context, creds *User) error
	SaveUser(ctx context.Context, creds *User) error
	PutUser(ctx context.Context, creds *User, userID string) error
	GetAllUser(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, userID string) (*User, error)
	GetUserByName(ctx context.Context, username string) (*User, error)
}
