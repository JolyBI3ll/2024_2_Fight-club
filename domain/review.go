package domain

type Review struct {
	ID     int     `gorm:"primary_key;auto_increment;column:id" json:"id"`
	UserID string  `gorm:"column:userId;not null" json:"userId"`
	HostID string  `gorm:"column:hostId;not null" json:"hostId"`
	Text   string  `gorm:"type:text;size:1000;column:text;not null" json:"text"`
	Rating float64 `gorm:"type:numeric;column:rating" json:"rating"`
	User   User    `gorm:"foreignkey:UserID;references:UUID"`
	Host   User    `gorm:"foreignkey:HostID;references:UUID"`
}
