package response

import (
	"github.com/gofiber/fiber/v2"
)

// APIResponse represents standard API response wrapper
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Success sends a successful response with optional data
func Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status: "success",
		Data:   data,
	})
}

// SuccessWithMessage sends a successful response with a message
func SuccessWithMessage(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "success",
		Message: message,
	})
}

// Error sends an error response with a message
func Error(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Status:  "fail",
		Message: message,
	})
}

// ValidationError sends a validation error response with field errors
func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Status: "fail",
		Errors: errors,
	})
}

// InternalError sends an internal server error response
func InternalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		Status:  "error",
		Message: message,
	})
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
		Status:  "fail",
		Message: message,
	})
}

// NotFound sends a not found response
func NotFound(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return c.Status(fiber.StatusNotFound).JSON(APIResponse{
		Status:  "fail",
		Message: message,
	})
}

// BadRequest sends a bad request response
func BadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Status:  "fail",
		Message: message,
	})
}

// Conflict sends a conflict response (e.g., duplicate entries)
func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(APIResponse{
		Status:  "fail",
		Message: message,
	})
}

// Created sends a created response with the created resource
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Status: "success",
		Data:   data,
	})
}

// OK sends an OK response with data
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status: "success",
		Data:   data,
	})
}

// NoContent sends a no content response (204)
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Forbidden sends a forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return c.Status(fiber.StatusForbidden).JSON(APIResponse{
		Status:  "fail",
		Message: message,
	})
}
