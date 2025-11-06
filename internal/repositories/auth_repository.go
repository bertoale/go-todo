package repositories

import (
	"rest-api/internal/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindEmailOrUsername(email, username string) (*models.User, error)
	Register(user *models.User) error
	FindByID(id uint) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

// FindByEmail implements AuthRepository.
func (a *authRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err:= a.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID implements AuthRepository.
func (a *authRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err:= a.db.First(&user, id).Error; err != nil {
		return nil, err
	}	
	return &user, nil
}

// FindEmailOrUsername implements AuthRepository.
func (a *authRepository) FindEmailOrUsername(email string, username string) (*models.User, error) {
	var user models.User
	if err:= a.db.Where("email = ? OR username = ?", email, username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Register implements AuthRepository.
func (a *authRepository) Register(user *models.User) error {
	return a.db.Create(user).Error
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
