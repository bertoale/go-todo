package main

import (
	"fmt"
	"log"
	"rest-api/config"
	"rest-api/internal/database"
	"rest-api/internal/middlewares"
	"rest-api/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
		BodyLimit: 10 * 1024 * 1024, // 10 MB
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CorsOrigin,
		AllowCredentials: true,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the REST API",
			"version": "1.0.0",
			"timestamp": fiber.Map{},
		})
	})

	routes.SetupRoutes(app, cfg)

	app.Use(middlewares.NotFound)

	port := cfg.Port
	log.Printf("🚀 Server is running on port %s", port)
	log.Printf("📍 Local: http://localhost:%s", port)
	log.Printf("🌍 Environment: %s", cfg.NodeEnv)

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ Unable to start server: %v", err)
	}


}