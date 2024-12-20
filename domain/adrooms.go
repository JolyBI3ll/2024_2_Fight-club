package domain

type AdRooms struct {
	ID           int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	AdID         string `gorm:"column:adId;not null" json:"adId"`
	Type         string `gorm:"type:text;size:100;column:type" json:"type"`
	SquareMeters int    `gorm:"column:squareMeters" json:"squareMeters"`
	Ad           Ad     `gorm:"foreignKey:adId;references:UUID"`
}

type AdRoomsResponse struct {
	Type         string `json:"type"`
	SquareMeters int    `json:"squareMeters"`
}
