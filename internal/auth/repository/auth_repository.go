package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(creds *domain.User) error {
	if err := r.db.Create(creds).Error; err != nil {
		return err
	}
	return nil
}

func (r *authRepository) PutUser(creds *domain.User, userID string) error {
	if err := r.db.Model(&domain.User{}).Where("UUID = ?", userID).Updates(creds).Error; err != nil {
		return err
	}
	return nil
}

func (r *authRepository) GetAllUser() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *authRepository) GetUserById(userID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("UUID = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByName(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
