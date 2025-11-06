// Package routes contains all route definitions for the blog API
// Routes dipisahkan per resource untuk kemudahan maintenance
package routes

import (
	"rest-api/config"
	"rest-api/internal/controllers"
	"rest-api/internal/database"
	"rest-api/internal/repositories"
	"rest-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes adalah function utama untuk setup semua routes aplikasi
// Menggunakan Dependency Injection pattern: Repository → Service → Controller
// Function ini dipanggil dari main.go setelah inisialisasi Fiber app
// Parameter:
//   - app: Fiber app instance
//   - cfg: Configuration object yang berisi environment variables
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize User Repository, Service, dan Controller dengan dependency injection
	userRepo := repositories.NewUserRepository(database.GetDB())
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)
	SetupUserRoutes(app, cfg, userController)
	authRepo := repositories.NewAuthRepository(database.GetDB())
	authService := services.NewAuthService(authRepo, cfg)
	authController := controllers.NewAuthController(authService, cfg)
	SetupAuthRoutes(app, cfg, authController)
}
