package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-fiber-jwt/pkg/mapper"
	"github.com/golang-fiber-jwt/pkg/response"
	"github.com/golang-fiber-jwt/pkg/validator"
)

// ParseAndValidate parses request body and validates it
func ParseAndValidate(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return response.ValidationError(c, errors)
	}
	return nil
}

// ParseValidateAndMap parses, validates, and auto-maps request to domain struct
// Returns the mapped domain struct and error
// Usage: data, err := handler.ParseValidateAndMap[RequestDTO, DomainStruct](c)
func ParseValidateAndMap[TReq any, TDomain any](c *fiber.Ctx) (*TDomain, error) {
	var req TReq

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		response.BadRequest(c, err.Error())
		return nil, err
	}

	// Validate request
	if errors := validator.ValidateStruct(req); errors != nil {
		response.ValidationError(c, errors)
		return nil, fiber.ErrBadRequest
	}

	// Auto-map to domain struct
	domain, err := mapper.AutoMap[TDomain](&req)
	if err != nil {
		response.InternalError(c, "Failed to process request")
		return nil, err
	}

	return domain, nil
}

// MapToDomain is a generic mapper from request to domain
// For cases where you already have the request parsed
func MapToDomain[TDomain any](req interface{}) (*TDomain, error) {
	return mapper.AutoMap[TDomain](req)
}
