package controllers

import (
	"fmt"
	"rest-api/internal/dto/request"
	"rest-api/internal/models"
	"rest-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var userID uint
	if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	userResponse, err := ctrl.userService.GetUserByID(userID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = fiber.StatusNotFound
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": userResponse,
	})
}

func (ctrl *UserController) UpdateUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	id := c.Params("id")

	var targetUserID uint
	if _, err := fmt.Sscanf(id, "%d", &targetUserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid user ID",
		})
	}

	

	var req request.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Call service untuk update user
	userResponse, err := ctrl.userService.UpdateUser(
		user.ID,
		targetUserID,
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "user tidak ditemukan" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "unauthorized" {
			statusCode = fiber.StatusForbidden
		} else if err.Error() == "email sudah digunakan" || err.Error() == "username sudah digunakan" {
			statusCode = fiber.StatusBadRequest
		}
		return c.Status(statusCode).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Profil berhasil diupdate.",
		"user":    userResponse,
	})
}

func (ctrl *UserController) GetProfile (c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	userResponse,  err := ctrl.userService.GetProfile(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": userResponse,
	})
}