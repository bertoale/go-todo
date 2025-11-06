package services

import (
	"errors"
	"rest-api/internal/dto/response"
	"rest-api/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(id uint) (*response.UserResponse, error)
	UpdateUser(currentUserID, targetUserID uint, username, email, password *string) (*response.UserResponse, error)
	CheckUsernameAvailability(username string, excludeUserID uint) error
	CheckEmailAvailability(email string, excludeUserID uint) error
	GetProfile(userID uint) (*response.UserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

// GetProfile implements UserService.
func (s *userService) GetProfile(userID uint) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, errors.New("gagal mengambil profil")
	}
	userResponse := &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResponse, nil
}

// CheckEmailAvailability implements UserService.
func (s *userService) CheckEmailAvailability(email string, excludeUserID uint) error {
	existingUser, err := s.userRepo.FindByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("gagal memeriksa email")
	}
	if existingUser != nil && existingUser.ID != excludeUserID {
		return errors.New("email sudah digunakan")
	}
	return nil
}

// CheckUsernameAvailability implements UserService.
func (s *userService) CheckUsernameAvailability(username string, excludeUserID uint) error {
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("gagal memeriksa username")
	}
	if existingUser != nil && existingUser.ID != excludeUserID {
		return errors.New("username sudah digunakan")
	}
	return nil
}

// GetUserByID implements UserService.
func (s *userService) GetUserByID(id uint) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user")
	}

	userResponse := &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResponse, nil
}

// UpdateUser implements UserService.
func (s *userService) UpdateUser(currentUserID uint, targetUserID uint, username *string, email *string, password *string) (*response.UserResponse, error) {
	user, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user")
	}
	if currentUserID != targetUserID {
		return nil, errors.New("unauthorized to update this user")
	}

	if email != nil && *email != user.Email {
		if err := s.CheckEmailAvailability(*email, currentUserID); err != nil {
			return nil, err
		}
		user.Email = *email
	}
	if username != nil && *username != user.Username {
		if err := s.CheckUsernameAvailability(*username, currentUserID); err != nil {
			return nil, err
		}
		user.Username = *username
	}

	if password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	userResponse := &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return userResponse, nil
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
