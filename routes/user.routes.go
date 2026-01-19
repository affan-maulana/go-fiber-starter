package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-fiber-jwt/internal/middleware"
	"github.com/golang-fiber-jwt/internal/user"
)

func UserRoutes(router fiber.Router, handler *user.Handler) {
	router.Route("/users", func(userRouter fiber.Router) {
		userRouter.Get("/", middleware.DeserializeUser, handler.ListUsers)
		userRouter.Get("/:id", middleware.DeserializeUser, handler.GetUserByID)
		// userRouter.Post("/", middleware.DeserializeUser, middleware.RequireAdminRole, handler.CreateUser)
		// userRouter.Put("/:id", middleware.DeserializeUser, middleware.RequireAdminRole, handler.UpdateUser)
		// userRouter.Delete("/:id", middleware.DeserializeUser, middleware.RequireAdminRole, handler.DeleteUser)
		// userRouter.Patch("/:id/restore", middleware.DeserializeUser, middleware.RequireAdminRole, handler.RestoreUser)
	})
}
