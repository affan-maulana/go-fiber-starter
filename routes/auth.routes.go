package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-fiber-jwt/internal/auth"
	"github.com/golang-fiber-jwt/internal/middleware"
)

func AuthRoutes(router fiber.Router, handler *auth.Handler) {
	router.Route("/auth", func(authRouter fiber.Router) {
		authRouter.Post("/register", handler.SignUpUser)
		authRouter.Post("/login", handler.SignInUser)
		authRouter.Get("/logout", middleware.DeserializeUser, handler.LogoutUser)
	})

	// User routes within auth domain
	router.Get("/user/me", middleware.DeserializeUser, handler.GetMe)
}
