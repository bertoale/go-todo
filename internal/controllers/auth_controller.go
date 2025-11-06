package controllers

import (
	"rest-api/config"
	"rest-api/internal/dto/request"
	"rest-api/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService
	cfg					*config.Config
}

func NewAuthController(authService services.AuthService, cfg *config.Config) *AuthController {
	return &AuthController{
		authService: authService,
		cfg: cfg,
	}
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req request.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body.",
		})
	}

	// Validasi input
	// if err := ctrl.validator.Struct(req); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"message": err.Error(),
	// 	})
	// }

	// Call service untuk login
	token, userResponse, err := ctrl.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Set cookie dengan token
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   ctrl.cfg.NodeEnv == "production",
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "Login successfully.",
		"token":   token,
		"user":    userResponse,
	})
}



func (ctrl *AuthController)Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	userResponse, err := ctrl.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    userResponse,
	})
}
	

