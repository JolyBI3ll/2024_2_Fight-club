package domain

import "time"

type Request struct {
	ID          int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	AdID        string    `gorm:"column:adId;not null" json:"adId"`
	UserID      string    `gorm:"column:userId;not null" json:"userId"`
	Status      string    `gorm:"column:status;not null;default:pending" json:"status"`
	CreatedDate time.Time `gorm:"type:timestamp;column:createDate;default:CURRENT_TIMESTAMP" json:"createDate"`
	UpdateDate  time.Time `gorm:"type:timestamp;column:updateDate" json:"updateDate"`
	CloseDate   time.Time `gorm:"type:timestamp;column:closeDate" json:"closeDate"`
	User        User      `gorm:"foreignkey:UserID;references:UUID"`
	Ad          Ad        `gorm:"foreignkey:AdID;references:UUID"`
}
