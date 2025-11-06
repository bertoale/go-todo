package routes

import (
	"rest-api/config"
	"rest-api/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, cfg *config.Config, authCtrl *controllers.AuthController) {
	users := app.Group("/api/auth")

	// GET /api/users/:id
	// Public route untuk mendapatkan informasi user by ID
	// Params: id (user ID)
	// Response: { user: { id, username, email, bio, avatar, createdAt } }
	users.Post("/register/", authCtrl.Register) 
	// PUT /api/users/:id
	// Protected route untuk update user profile
	// Headers: Authorization: Bearer <token>
	// Params: id (user ID)
	// Request body (multipart/form-data): { username?, email?, bio?, password?, avatar? }
	// Response: { message, user }
	// Note: User hanya bisa update profile sendiri
	users.Post("/login/", authCtrl.Login)
}
