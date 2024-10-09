package ds

type User struct {
	UUID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username   string
	Password   string
	Email      string
	Name       string
	Score      float32
	Avatar     string
	Sex        rune
	GuestCount int
	Age        int
	Address    string
	IsHost     bool
	Ads        []Ad      `gorm:"foreignKey:AuthorUUID"`
	Requests   []Request `gorm:"foreignKey:UserID"`
	Reviews    []Review  `gorm:"foreignKey:UserID"`
}

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

type Request struct {
	ID             string `gorm:"primaryKey"`
	AdID           string
	UserID         string
	Status         string
	RequestDate    string
	RequestedDates []string `gorm:"type:text[]"`
}

type Review struct {
	ID     string `gorm:"primaryKey"`
	UserID string
	HostId string
	Text   string
	Rating float32
}
