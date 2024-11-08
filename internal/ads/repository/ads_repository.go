package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type adRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) domain.AdRepository {
	return &adRepository{
		db: db,
	}
}

func (r *adRepository) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAllPlaces called", zap.String("request_id", requestID))
	var ads []domain.GetAllAdsResponse

	query := r.db.Model(&domain.Ad{}).Joins("JOIN cities ON  ads.\"cityId\" = cities.id").Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").
		Select("ads.*, cities.title as cityName")

	if filter.Location != "" {
		query = query.Where("cities.\"enTitle\" = ?", filter.Location)
	}

	if filter.Rating != "" {
		rating, err := strconv.ParseFloat(filter.Rating, 1)
		if err != nil {
			logger.DBLogger.Error("Invalid rating value", zap.String("request_id", requestID))
			return nil, errors.New("Invalid rating value")
		}
		query = query.Where("users.score >= ?", rating)
	}

	if filter.NewThisWeek == "true" {
		lastWeek := time.Now().AddDate(0, 0, -7)
		query = query.Where("\"publicationDate\" >= ?", lastWeek)
	}

	if filter.HostGender != "" && filter.HostGender != "any" {
		if filter.HostGender == "male" {
			query = query.Where("users.sex = ?", "M")
		} else if filter.HostGender == "female" {
			query = query.Where("users.sex = ?", "F")
		}
	}

	if filter.GuestCount != "" {
		switch filter.GuestCount {
		case "5":
			query = query.Where("users.\"guestCount\" > ?", 5)
		case "10":
			query = query.Where("users.\"guestCount\" > ?", 10)
		case "20":
			query = query.Where("users.\"guestCount\" > ?", 20)
		case "50":
			query = query.Where("users.\"guestCount\" > ?", 50)
		}
	}

	if filter.Offset != 0 {
		query = query.Offset(filter.Offset)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching all places", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("Error fetching all places")
	}

	for i, ad := range ads {
		var images []domain.Image
		var user domain.User
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("Error fetching images for ad")
		}

		err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("Error fetching user")
		}
		ads[i].AdAuthor.Name = user.Name
		ads[i].AdAuthor.Avatar = user.Avatar
		ads[i].AdAuthor.Rating = user.Score

		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched all places", zap.String("request_id", requestID), zap.Int("count", len(ads)))
	return ads, nil
}

func (r *adRepository) GetPlaceById(ctx context.Context, adId string) (domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPlaceById called", zap.String("adId", adId), zap.String("request_id", requestID))

	var ad domain.GetAllAdsResponse

	query := r.db.Model(&domain.Ad{}).Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").Joins("JOIN cities ON  ads.\"cityId\" = cities.id").
		Select("ads.*, cities.title as cityName").Where("ads.uuid = ?", adId)

	if err := query.Find(&ad).Error; err != nil {
		logger.DBLogger.Error("Error fetching place", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("Error fetching place")
	}

	var images []domain.Image
	var user domain.User
	err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
	if err != nil {
		logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("Error fetching images for ad")
	}

	err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
	if err != nil {
		logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("Error fetching user")
	}

	ad.AdAuthor.Name = user.Name
	ad.AdAuthor.Avatar = user.Avatar
	ad.AdAuthor.Rating = user.Score

	for _, img := range images {
		ad.Images = append(ad.Images, domain.ImageResponse{
			ID:        img.ID,
			ImagePath: img.ImageUrl,
		})
	}

	logger.DBLogger.Info("Successfully fetched place by ID", zap.String("adId", adId), zap.String("request_id", requestID))
	return ad, nil
}

func (r *adRepository) CreatePlace(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreatePlace called", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	var city domain.City
	if err := r.db.Where("title = ?", newAd.CityName).First(&city).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("Error finding city")
	}
	ad.CityID = city.ID
	ad.PublicationDate = time.Now()
	if err := r.db.Create(ad).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("Error creating place")
	}

	logger.DBLogger.Info("Successfully place", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) SavePlace(ctx context.Context, ad *domain.Ad) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("SavePlace called", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	if err := r.db.Save(ad).Error; err != nil {
		logger.DBLogger.Error("Error saving place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return err
	}
	logger.DBLogger.Info("Successfully saved place", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) UpdatePlace(ctx context.Context, ad *domain.Ad, adId string, userId string, updatedPlace domain.UpdateAdRequest) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("UpdatePlace called", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))

	var oldAd domain.Ad
	if err := r.db.Where("uuid = ?", adId).First(&oldAd).Error; err != nil {
		logger.DBLogger.Error("Ad not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad not found")
	}

	if oldAd.AuthorUUID != userId {
		logger.DBLogger.Warn("User is not the owner of the ad", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("not owner of ad")
	}
	var city domain.City
	if err := r.db.Where("title = ?", updatedPlace.CityName).First(&city).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("Error finding city")
	}
	ad.CityID = city.ID
	if err := r.db.Model(&oldAd).Updates(ad).Error; err != nil {
		logger.DBLogger.Error("Error updating place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully updated place", zap.String("adId", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) DeletePlace(ctx context.Context, adId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeletePlace called", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))

	var ad domain.Ad
	if err := r.db.Where("uuid = ?", adId).First(&ad).Error; err != nil {
		logger.DBLogger.Error("Ad not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad not found")
	}

	if ad.AuthorUUID != userId {
		logger.DBLogger.Warn("User is not the owner of the ad", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("not owner of ad")
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.Image{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.AdPosition{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.AdAvailableDate{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.Request{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	if err := r.db.Delete(&ad).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully deleted place", zap.String("adId", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPlacesPerCity called", zap.String("city", city), zap.String("request_id", requestID))

	var ads []domain.GetAllAdsResponse
	query := r.db.Model(&domain.Ad{}).Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").Joins("JOIN cities ON  ads.\"cityId\" = cities.id").
		Select("ads.*, cities.title as cityName").Where("cities.\"enTitle\" = ?", city)
	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching places per city", zap.String("city", city), zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	for i, ad := range ads {
		var images []domain.Image
		var user domain.User
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("Error fetching images for ad")
		}

		err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("Error fetching user")
		}
		ads[i].AdAuthor.Name = user.Name
		ads[i].AdAuthor.Avatar = user.Avatar
		ads[i].AdAuthor.Rating = user.Score

		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched places per city", zap.String("city", city), zap.Int("count", len(ads)), zap.String("request_id", requestID))
	return ads, nil
}

func (r *adRepository) SaveImages(ctx context.Context, adUUID string, imagePaths []string) error {
	for _, path := range imagePaths {
		image := domain.Image{
			AdID:     adUUID,
			ImageUrl: path,
		}
		if err := r.db.Create(&image).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *adRepository) GetAdImages(ctx context.Context, adId string) ([]string, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAdImages called", zap.String("request_id", requestID), zap.String("adId", adId))

	var imageUrls []string

	err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", adId).Pluck("imageUrl", &imageUrls).Error
	if err != nil {
		logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("Error fetching images for ad")
	}

	logger.DBLogger.Info("Successfully fetched images for ad", zap.String("request_id", requestID), zap.Int("count", len(imageUrls)))
	return imageUrls, nil
}

func (r *adRepository) GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserPlaces called", zap.String("city", userId), zap.String("request_id", requestID))

	var ads []domain.GetAllAdsResponse
	query := r.db.Model(&domain.Ad{}).Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").Joins("JOIN cities ON  ads.\"cityId\" = cities.id").
		Select("ads.*, users.avatar, users.name, users.score as rating , cities.title as cityName").Where("users.uuid = ?", userId)
	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching user places", zap.String("city", userId), zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	for i, ad := range ads {
		var images []domain.Image
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("Error fetching images for ad")
		}

		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched user places", zap.String("city", userId), zap.Int("count", len(ads)), zap.String("request_id", requestID))
	return ads, nil
}

func (r *adRepository) DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) (string, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeleteAdImage called", zap.String("ad", adId), zap.Int("image", imageId), zap.String("request_id", requestID))

	var ad domain.Ad
	if err := r.db.First(&ad, "uuid = ?", adId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("ad not found")
		}
		return "", errors.New("Error fetching ad")
	}

	if ad.AuthorUUID != userId {
		return "", errors.New("not owner of ad")
	}

	var image domain.Image
	if err := r.db.First(&image, "id = ? AND \"adId\" = ?", imageId, adId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("image not found")
		}
		return "", errors.New("error finding image")
	}

	if err := r.db.Delete(&image).Error; err != nil {
		return "", errors.New("error deleting image from database")
	}

	logger.DBLogger.Info("Image deleted successfully", zap.Int("image_id", imageId), zap.String("ad_id", adId), zap.String("request_id", requestID))
	return image.ImageUrl, nil
}
