package domain

type Request struct {
	ID             string `gorm:"primaryKey"`
	AdID           string
	UserID         string
	Status         string
	RequestDate    string
	RequestedDates []string `gorm:"type:text[]"`
}
