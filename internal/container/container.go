package container

import (
	"github.com/golang-fiber-jwt/internal/auth"
	"github.com/golang-fiber-jwt/internal/user"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	AuthHandler *auth.Handler
	UserHandler *user.Handler
	// Add other handlers here as you create new modules
	// ProductHandler *product.Handler
	// OrderHandler   *order.Handler
}

// NewContainer creates a new dependency injection container
func NewContainer(db *gorm.DB) *Container {
	// Auth
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// User
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Wire other modules here
	// productRepo := postgresql.NewProductRepository(db)
	// productService := product.NewAuthService(productRepo)
	// productHandler := product.NewAuthHandler(productService)

	return &Container{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		// ProductHandler: productHandler,
	}
}
