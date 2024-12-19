package domain

//go:generate easyjson -all regions.go

import (
	"context"
	"time"
)

//easyjson:json
type VisitedRegionsList []VisitedRegions

type VisitedRegions struct {
	ID             int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	Name           string    `gorm:"text;size:1000;column:name" json:"name"`
	UserID         string    `gorm:"column:userId;not null" json:"userId"`
	StartVisitDate time.Time `gorm:"type:timestamp;column:startVisitDate;default:CURRENT_TIMESTAMP" json:"startVisitDate"`
	EndVisitDate   time.Time `gorm:"type:timestamp;column:endVisitDate;default:CURRENT_TIMESTAMP" json:"endVisitDate"`
	User           User      `gorm:"foreignKey:UserID;references:UUID" json:"-"`
}

type RegionRepository interface {
	GetVisitedRegions(ctx context.Context, userId string) ([]VisitedRegions, error)
}
