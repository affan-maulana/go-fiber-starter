package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-fiber-jwt/config"
	"github.com/golang-fiber-jwt/internal/container"
	"github.com/golang-fiber-jwt/routes"
)

func init() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	config.ConnectDB(&cfg)
}

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3334",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
		AllowCredentials: true,
	}))

	// Initialize dependency injection container
	c := container.NewContainer(config.DB)

	// Setup routes with injected handlers
	routes.SetupRoutes(app, c.AuthHandler, c.UserHandler)

	log.Fatal(app.Listen(":3334"))
}
