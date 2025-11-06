package services

import (
	"errors"
	"rest-api/config"
	"rest-api/internal/dto/response"
	"rest-api/internal/middlewares"
	"rest-api/internal/models"
	"rest-api/internal/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(username, email, password string) (*response.UserResponse, error)
	Login(email, password string) (string, *response.UserResponse, error)
	GenerateToken(userID uint) (string, error)
}

type authService struct {
	authRepo repositories.AuthRepository
	cfg      *config.Config
}

// GenerateToken implements AuthService.
func (a *authService) GenerateToken(userID uint) (string, error) {
	return middlewares.GenerateToken(userID, a.cfg)
}


// Login implements AuthService.
func (a *authService) Login(email string, password string) (string, *response.UserResponse, error) {
	if email == "" || password == "" {
		return "", nil, errors.New("email and password are required")
	}

	user, err := a.authRepo.FindByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, errors.New("email not found")
		}
		return "", nil, errors.New("failed to retrieve user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password));  err != nil {
		return "", nil, errors.New("incorrect password")
	}

	token, err := a.GenerateToken(user.ID)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	userResponse := &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return token, userResponse, nil
}


// Register implements AuthService.
func (a *authService) Register(username string, email string, password string)(*response.UserResponse, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	existingUser,err := a.authRepo.FindEmailOrUsername(email, username)
	if err == nil && existingUser != nil {
		return nil, errors.New("email or username already in use")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("failed to check existing users")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := a.authRepo.Register(user); err != nil {
		return nil, errors.New("failed to register user")
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



func NewAuthService(authRepo repositories.AuthRepository, cfg *config.Config) AuthService {
	return &authService{authRepo: authRepo, cfg: cfg}
}

func (s *authService) GetTokenExpiration() time.Duration {
	return 7 * 24 * time.Hour
}
