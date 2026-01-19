package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-fiber-jwt/config"
	"github.com/golang-fiber-jwt/internal/user"
	"github.com/golang-fiber-jwt/pkg/handler"
	"github.com/golang-fiber-jwt/pkg/response"
	"github.com/golang-jwt/jwt"
)

// Handler handles HTTP requests for auth domain
// This layer is allowed to import Fiber for HTTP handling
type Handler struct {
	service Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Helper function to handle service errors with appropriate HTTP status codes
func (h *Handler) handleServiceError(c *fiber.Ctx, err error) error {
	errorMessage := err.Error()

	switch errorMessage {
	case "user with that email already exists":
		return response.Conflict(c, errorMessage)
	case "passwords do not match", "name is required", "email is required", "password is required", "password must be at least 8 characters":
		return response.BadRequest(c, errorMessage)
	case "invalid email or password":
		return response.BadRequest(c, errorMessage)
	case "user not found":
		return response.NotFound(c, errorMessage)
	default:
		return response.Error(c, fiber.StatusBadGateway, errorMessage)
	}
}

// Helper function to map domain User to transport UserResponse
func (h *Handler) userToResponse(user *user.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Photo:     user.Photo,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// SignUpUser handles user registration requests
func (h *Handler) SignUpUser(c *fiber.Ctx) error {
	// Parse, validate, and map request to domain
	signUpData, err := handler.ParseValidateAndMap[SignUpRequest, SignUpData](c)
	if err != nil {
		return nil // Response already sent by helper
	}

	// Call service
	user, err := h.service.SignUp(signUpData)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Map to response and return success
	userResponse := h.userToResponse(user)
	return response.Created(c, UserDataResponse{
		User: userResponse,
	})
}

// SignInUser handles user login requests
func (h *Handler) SignInUser(c *fiber.Ctx) error {
	var req SignInRequest

	// Parse and validate request
	if err := handler.ParseAndValidate(c, &req); err != nil {
		return err
	}

	// Call service
	_, user, err := h.service.SignIn(req.Email, req.Password)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Generate JWT token (technical concern - stays in handler)
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return response.InternalError(c, "Failed to load config")
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID.String()
	claims["exp"] = now.Add(cfg.JwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return response.InternalError(c, "Failed to generate token")
	}

	// Set cookie (HTTP concern - stays in handler)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   cfg.JwtMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(AuthResponse{
		Status: "success",
		Token:  tokenString,
	})
}

// LogoutUser handles user logout requests
func (h *Handler) LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return response.SuccessWithMessage(c, fiber.StatusOK, "Logged out successfully")
}

// GetMe returns the current authenticated user
func (h *Handler) GetMe(c *fiber.Ctx) error {
	// Get user ID from context (set by middleware)
	userID := c.Locals("userId")
	if userID == nil {
		return response.Unauthorized(c, "Unauthorized")
	}

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Map to response
	userResponse := h.userToResponse(user)
	return response.OK(c, UserDataResponse{
		User: userResponse,
	})
}
