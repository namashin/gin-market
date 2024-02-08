package repository

import (
	"errors"
	"gin-market/models"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(user models.User) error
	FindUser(email string) (*models.User, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) CreateUser(user models.User) error {
	result := ar.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (ar *AuthRepository) FindUser(email string) (*models.User, error) {
	var user models.User

	result := ar.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}
