package domain

type Image struct {
	ID       int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	AdID     string `gorm:"column:adId;not null" json:"adId"`
	ImageUrl string `gorm:"type:text;size:1000;column:imageUrl" json:"imageUrl"`
	Ad       Ad     `gorm:"foreignKey:adId;references:UUID"`
}

type ImageResponse struct {
	ID        int    `json:"id"`
	ImagePath string `json:"path"`
}
