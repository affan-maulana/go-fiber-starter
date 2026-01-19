package auth

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents user data for HTTP responses
type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse represents authentication response with token
type AuthResponse struct {
	Status string            `json:"status"`
	Token  string            `json:"token,omitempty"`
	Data   *UserDataResponse `json:"data,omitempty"`
}

// UserDataResponse wraps user data for responses
type UserDataResponse struct {
	User UserResponse `json:"user"`
}

type SignUpRequest struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=8"`
	Photo           string `json:"photo"`
}

// SignInRequest represents user login HTTP request
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
