package routes

import (
	"rest-api/config"
	"rest-api/internal/controllers"
	"rest-api/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(app *fiber.App, cfg *config.Config, taskCtrl *controllers.TaskController) {
	tasks := app.Group("/api/tasks")
	tasks.Get("/:id", middlewares.Auth(cfg), taskCtrl.GetTaskByID)
	tasks.Get("/", middlewares.Auth(cfg), taskCtrl.GetTasksByUserID)
	tasks.Post("/", middlewares.Auth(cfg), taskCtrl.CreateTask)
	tasks.Put("/:id", middlewares.Auth(cfg), taskCtrl.UpdateTask)
	tasks.Delete("/:id", middlewares.Auth(cfg), taskCtrl.DeleteTask)
}
