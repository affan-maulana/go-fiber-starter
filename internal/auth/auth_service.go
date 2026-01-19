package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-fiber-jwt/internal/user"
	"github.com/golang-fiber-jwt/pkg/hashing"
	"github.com/google/uuid"
)

// Service defines the interface for auth business logic
type Service interface {
	SignUp(data *SignUpData) (*user.User, error)
	SignIn(email, password string) (token string, user *user.User, err error)
	GetUserByID(id string) (*user.User, error)
}

// service implements the Service interface
// Pure business logic - no framework dependencies
type service struct {
	repo Repository
}

// NewAuthService creates a new auth service
func NewAuthService(repo Repository) Service {
	return &service{repo: repo}
}

// SignUp handles user registration business logic
func (s *service) SignUp(data *SignUpData) (*user.User, error) {
	// Validate input
	if err := s.validateSignUpData(data); err != nil {
		return nil, err
	}

	// Check if passwords match
	if data.Password != data.PasswordConfirm {
		return nil, fmt.Errorf("passwords do not match")
	}

	// Hash password
	hashedPassword, err := hashing.HashPassword(data.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	now := time.Now()
	user := &user.User{
		ID:        uuid.New(),
		Name:      data.Name,
		Email:     strings.ToLower(strings.TrimSpace(data.Email)),
		Password:  hashedPassword,
		Role:      "user",
		Provider:  "local",
		Photo:     s.getPhotoOrDefault(data.Photo),
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to repository
	if err := s.repo.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("user with that email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// SignIn handles user authentication business logic
func (s *service) SignIn(email, password string) (string, *user.User, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := hashing.VerifyPassword(user.Password, password); err != nil {
		return "", nil, fmt.Errorf("invalid email or password")
	}

	// Generate token (simple random token for now - will be enhanced in handler layer with JWT)
	token, err := hashing.GenerateToken()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetUserByID retrieves a user by their ID
func (s *service) GetUserByID(id string) (*user.User, error) {
	return s.repo.GetUserByID(id)
}

// validateSignUpData validates sign up data
func (s *service) validateSignUpData(data *SignUpData) error {
	if data.Name == "" {
		return fmt.Errorf("name is required")
	}
	if data.Email == "" {
		return fmt.Errorf("email is required")
	}
	if data.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(data.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// getPhotoOrDefault returns the photo or default photo
func (s *service) getPhotoOrDefault(photo string) string {
	if photo == "" {
		return "default.png"
	}
	return photo
}
