package routes

import (
	"rest-api/config"
	"rest-api/internal/controllers"
	"rest-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, cfg *config.Config, userCtrl *controllers.UserController) {
	users := app.Group("/api/users")
	users.Put("/:id", middlewares.Auth(cfg), userCtrl.UpdateUser)
	users.Get("/", middlewares.Auth(cfg), userCtrl.GetProfile)

}
