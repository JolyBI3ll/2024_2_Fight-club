package domain

import "time"

type AdAvailableDate struct {
	ID                int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	AdID              string    `gorm:"column:adId;not null" json:"adId"`
	AvailableDateFrom time.Time `gorm:"type:date;column:availableDateFrom" json:"availableDateFrom"`
	AvailableDateTo   time.Time `gorm:"type:date;column:availableDateTo" json:"availableDateTo"`
	Ad                Ad        `gorm:"foreignKey:adId;references:UUID"`
}
