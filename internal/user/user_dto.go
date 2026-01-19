package user

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents user creation HTTP request
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
	Provider string `json:"provider" validate:"omitempty,oneof=local google facebook"`
	Photo    string `json:"photo"`
	Verified bool   `json:"verified"`
}

// UpdateUserRequest represents user update HTTP request
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"omitempty,oneof=user admin"`
	Photo    string `json:"photo"`
	Verified bool   `json:"verified"`
}

// UserResponse represents user data for HTTP responses
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Provider  string     `json:"provider"`
	Photo     string     `json:"photo"`
	Verified  bool       `json:"verified"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

// UserDataResponse wraps user data for single user responses
type UserDataResponse struct {
	User UserResponse `json:"user"`
}

// ListUsersQuery represents query parameters for listing users
type ListUsersQuery struct {
	Page        int    `query:"page" validate:"omitempty,min=1"`
	PerPage     int    `query:"per_page" validate:"omitempty,min=1,max=100"`
	Search      string `query:"search"`
	SearchBy    string `query:"search_by" validate:"omitempty,oneof=name email"`
	Role        string `query:"role" validate:"omitempty,oneof=user admin"`
	Provider    string `query:"provider" validate:"omitempty,oneof=local google facebook"`
	Verified    *bool  `query:"verified"`
	ShowDeleted bool   `query:"show_deleted"`
}
