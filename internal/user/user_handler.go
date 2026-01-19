package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-fiber-jwt/pkg/handler"
	"github.com/golang-fiber-jwt/pkg/response"
)

// Handler handles HTTP requests for user domain
type Handler struct {
	service Service
}

// NewHandler creates a new user handler
func NewUserHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// handleServiceError maps service errors to appropriate HTTP responses
func (h *Handler) handleServiceError(c *fiber.Ctx, err error) error {
	errorMessage := err.Error()

	switch errorMessage {
	case "user with that email already exists", "email is already taken by another user":
		return response.Conflict(c, errorMessage)
	case "name is required", "email is required", "password is required", "password must be at least 8 characters":
		return response.BadRequest(c, errorMessage)
	case "invalid user ID format":
		return response.BadRequest(c, errorMessage)
	case "user not found":
		return response.NotFound(c, errorMessage)
	case "user is not deleted":
		return response.BadRequest(c, errorMessage)
	default:
		return response.InternalError(c, "Internal server error")
	}
}

// userToResponse maps domain UserResponse to UserResponse DTO
func (h *Handler) userToResponse(user *UserResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Provider:  user.Provider,
		Photo:     user.Photo,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

// ListUsers handles GET /users - retrieve users with pagination and filtering
func (h *Handler) ListUsers(c *fiber.Ctx) error {
	// Parse query parameters manually
	query := ListUsersQuery{
		Page:    1,
		PerPage: 10,
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}

	// Parse per_page
	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
			query.PerPage = perPage
		}
	}

	// Parse other parameters
	query.Search = c.Query("search")
	// query.Role = c.Query("role")
	// query.Provider = c.Query("provider")

	// Parse verified
	// if verifiedStr := c.Query("verified"); verifiedStr != "" {
	// 	if verified, err := strconv.ParseBool(verifiedStr); err == nil {
	// 		query.Verified = &verified
	// 	}
	// }

	// Parse show_deleted
	if showDeletedStr := c.Query("show_deleted"); showDeletedStr != "" {
		if showDeleted, err := strconv.ParseBool(showDeletedStr); err == nil {
			query.ShowDeleted = showDeleted
		}
	}

	// Use goroutines for concurrent processing
	type result struct {
		users      []UserResponse
		total      int64
		totalPages int
		err        error
	}

	resultChan := make(chan result, 1)

	// Start goroutine for data fetching and pagination calculation
	go func() {
		defer close(resultChan)

		// Fetch users and total count
		users, total, err := h.service.GetUsers(query)
		if err != nil {
			resultChan <- result{err: err}
			return
		}

		// Calculate pagination concurrently
		totalPages := h.service.CalculatePagination(total, query.Page, query.PerPage)

		resultChan <- result{
			users:      users,
			total:      total,
			totalPages: totalPages,
			err:        nil,
		}
	}()

	// Wait for result
	res := <-resultChan
	if res.err != nil {
		return h.handleServiceError(c, res.err)
	}

	// Return response
	return response.OK(c, UserListResponse{
		Items:      res.users,
		Total:      res.total,
		Page:       query.Page,
		PerPage:    query.PerPage,
		TotalPages: res.totalPages,
	})
}

// GetUserByID handles GET /users/:id - retrieve user by ID
func (h *Handler) GetUserByID(c *fiber.Ctx) error {
	// Get ID from URL parameters
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "user ID is required")
	}

	// Call service
	user, err := h.service.GetUserByID(id)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Map to response DTO
	userResponse := h.userToResponse(user)
	return response.OK(c, UserDataResponse{
		User: userResponse,
	})
}

// CreateUser handles POST /users - create new user
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	// Parse and validate request
	var req CreateUserRequest
	if err := handler.ParseAndValidate(c, &req); err != nil {
		return err
	}

	// Map to domain struct
	createData := &CreateUserData{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Provider: req.Provider,
		Photo:    req.Photo,
		Verified: req.Verified,
	}

	// Call service
	err := h.service.CreateUser(createData)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Map to response DTO and return success
	// userResponse := h.userToResponse(user)
	return response.Created(c, nil)
}

// UpdateUser handles PUT /users/:id - update existing user
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	// Get ID from URL parameters
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "user ID is required")
	}

	// Parse and validate request
	var req UpdateUserRequest
	if err := handler.ParseAndValidate(c, &req); err != nil {
		return err
	}

	// Map to domain struct
	updateData := &UpdateUserData{
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Photo:    req.Photo,
		Verified: req.Verified,
	}

	// Call service
	err := h.service.UpdateUser(id, updateData)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	return response.OK(c, nil)
}

// DeleteUser handles DELETE /users/:id - soft delete user
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	// Get ID from URL parameters
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "user ID is required")
	}

	// Call service
	if err := h.service.DeleteUser(id); err != nil {
		return h.handleServiceError(c, err)
	}

	return response.SuccessWithMessage(c, fiber.StatusOK, "User deleted successfully")
}

// RestoreUser handles POST /users/:id/restore - restore soft deleted user
func (h *Handler) RestoreUser(c *fiber.Ctx) error {
	// Get ID from URL parameters
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "user ID is required")
	}

	// Call service
	user, err := h.service.RestoreUser(id)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Map to response DTO and return success
	userResponse := h.userToResponse(user)
	return response.OK(c, UserDataResponse{
		User: userResponse,
	})
}

// GetUserStats handles GET /users/stats - get user statistics (bonus endpoint)
func (h *Handler) GetUserStats(c *fiber.Ctx) error {
	// Get total users
	totalQuery := ListUsersQuery{Page: 1, PerPage: 1}
	_, total, err := h.service.GetUsers(totalQuery)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Get verified users
	verified := true
	verifiedQuery := ListUsersQuery{Page: 1, PerPage: 1, Verified: &verified}
	_, totalVerified, err := h.service.GetUsers(verifiedQuery)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	// Get admin users
	adminQuery := ListUsersQuery{Page: 1, PerPage: 1, Role: "admin"}
	_, totalAdmins, err := h.service.GetUsers(adminQuery)
	if err != nil {
		return h.handleServiceError(c, err)
	}

	stats := map[string]interface{}{
		"total_users":      total,
		"verified_users":   totalVerified,
		"admin_users":      totalAdmins,
		"regular_users":    total - totalAdmins,
		"unverified_users": total - totalVerified,
	}

	return response.OK(c, stats)
}
