package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"gorm.io/gorm"
)

type adRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) domain.AdRepository {
	return &adRepository{
		db: db,
	}
}

func (r *adRepository) GetAllPlaces() ([]domain.Ad, error) {
	var ads []domain.Ad
	if err := r.db.Find(&ads).Error; err != nil {
		return nil, err
	}
	return ads, nil
}
