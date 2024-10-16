package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"errors"
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

func (r *adRepository) GetPlaceById(adId string) (domain.Ad, error) {
	var ad domain.Ad
	if err := r.db.Where("id = ?", adId).First(&ad).Error; err != nil {
		return ad, err
	}
	return ad, nil
}

func (r *adRepository) CreatePlace(ad *domain.Ad) error {
	if err := r.db.Create(ad).Error; err != nil {
		return err
	}
	if err := r.SavePlace(ad); err != nil {
		return err
	}
	return nil
}

func (r *adRepository) SavePlace(ad *domain.Ad) error {
	if err := r.db.Save(ad).Error; err != nil {
		return err
	}
	return nil
}

func (r *adRepository) UpdatePlace(ad *domain.Ad, adId string, userId string) error {
	var oldAd domain.Ad
	if err := r.db.Where("id = ?", adId).First(&oldAd).Error; err != nil {
		return errors.New("ad not found")
	}
	if oldAd.AuthorUUID != userId {
		return errors.New("not owner of ad")
	}

	if err := r.db.Model(&oldAd).Updates(ad).Error; err != nil {
		return err
	}
	return nil
}

func (r *adRepository) DeletePlace(adId string, userId string) error {
	var ad domain.Ad
	if err := r.db.Where("id = ?", adId).First(&ad).Error; err != nil {
		return errors.New("ad not found")
	}
	if ad.AuthorUUID != userId {
		return errors.New("not owner of ad")
	}
	if err := r.db.Delete(&ad).Error; err != nil {
		return err
	}
	return nil
}

func (r *adRepository) GetPlacesPerCity(city string) ([]domain.Ad, error) {
	var ads []domain.Ad
	if err := r.db.Where("location_main = ?", city).Find(&ads).Error; err != nil {
		return nil, err
	}
	return ads, nil
}
