package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-fiber-jwt/internal/auth"
	"github.com/golang-fiber-jwt/internal/user"
)

func SetupRoutes(app *fiber.App, authHandler *auth.Handler, userHandler *user.Handler) {
	micro := fiber.New()
	app.Mount("/api", micro)

	// Setup all module routes
	AuthRoutes(micro, authHandler)
	UserRoutes(micro, userHandler)

	// Health check
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Server is running! JWT Authentication with Golang, Fiber, and GORM",
		})
	})

	// 404 handler
	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})
}
