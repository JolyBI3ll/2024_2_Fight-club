package domain

import "time"

type AdAvailableDate struct {
	ID            int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	AdID          string    `gorm:"column:adId;not null" json:"adId"`
	AvailableDate time.Time `gorm:"type:date;column:availableDate" json:"availableDate"`
	Ad            Ad        `gorm:"foreignkey:AdId;references:UUID"`
}
