package domain

type AdPosition struct {
	ID        int     `gorm:"primary_key;auto_increment;column:id" json:"id"`
	AdID      string  `gorm:"column:adId;not null" json:"adId"`
	Latitude  float64 `gorm:"type:numeric;column:latitude" json:"latitude"`
	Longitude float64 `gorm:"type:numeric;column:longitude" json:"longitude"`
	Ad        Ad      `gorm:"foreignkey:AdId;references:UUID"`
}
