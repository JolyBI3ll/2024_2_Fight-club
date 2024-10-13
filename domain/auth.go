package domain

import (
	"time"
)

type User struct {
	UUID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username   string
	Password   string
	Email      string
	Name       string
	Score      float32
	Avatar     string
	Sex        rune
	GuestCount int
	Birthdate  time.Time `gorm:"type:timestamptz"`
	Address    string
	IsHost     bool
	Ads        []Ad      `gorm:"foreignKey:AuthorUUID"`
	Requests   []Request `gorm:"foreignKey:UserID"`
	Reviews    []Review  `gorm:"foreignKey:UserID"`
}

type AuthRepository interface {
	CreateUser(creds *User) error
	PutUser(creds *User, userID string) error
	GetAllUser() ([]User, error)
	GetUserById(userID string) (*User, error)
	GetUserByName(username string) (*User, error)
}
