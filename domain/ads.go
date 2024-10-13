package domain

type Ad struct {
	ID              string `gorm:"primaryKey"`
	LocationMain    string
	LocationStreet  string
	Position        []float64 `gorm:"type:float[]"`
	Images          []string  `gorm:"type:text[]"`
	AuthorUUID      string
	PublicationDate string
	AvailableDates  []string `gorm:"type:text[]"`
	Distance        float32
	Requests        []Request `gorm:"foreignKey:AdID"`
}
