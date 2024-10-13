package domain

type Review struct {
	ID     string `gorm:"primaryKey"`
	UserID string
	HostId string
	Text   string
	Rating float32
}
