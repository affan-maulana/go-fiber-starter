package user

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service defines the interface for user business logic
type Service interface {
	// GetUsers retrieves users with filtering and pagination
	GetUsers(query ListUsersQuery) ([]UserResponse, int64, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(id string) (*UserResponse, error)

	// CreateUser creates a new user in the system
	CreateUser(data *CreateUserData) error

	// UpdateUser updates an existing user
	UpdateUser(id string, data *UpdateUserData) error

	// DeleteUser soft deletes a user
	DeleteUser(id string) error

	// RestoreUser restores a soft deleted user
	RestoreUser(id string) (*UserResponse, error)

	// CalculatePagination calculates total pages for pagination
	CalculatePagination(total int64, page, perPage int) int
}

// service implements Service interface with pure business logic
type service struct {
	repo Repository
}

// NewUserService creates a new user service
func NewUserService(repo Repository) Service {
	return &service{repo: repo}
}

// GetUsers retrieves users with filtering and pagination
func (s *service) GetUsers(query ListUsersQuery) ([]UserResponse, int64, error) {
	// Business rule: Default pagination values
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PerPage <= 0 || query.PerPage > 100 {
		query.PerPage = 10
	}

	return s.repo.GetUsers(query)
}

// GetUserByID retrieves a user by their ID
func (s *service) GetUserByID(id string) (*UserResponse, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := s.repo.GetUserByID(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// CreateUser creates a new user in the system
func (s *service) CreateUser(data *CreateUserData) error {
	// Business rule validations
	if data.Name == "" {
		return errors.New("name is required")
	}

	if data.Email == "" {
		return errors.New("email is required")
	}

	if data.Password == "" {
		return errors.New("password is required")
	}

	if len(data.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(data.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return errors.New("user with that email already exists")
	}

	// Set default values
	if data.Role == "" {
		data.Role = "user"
	}

	if data.Provider == "" {
		data.Provider = "local"
	}

	if data.Photo == "" {
		data.Photo = "default.png"
	}

	// Create user entity
	user := &User{
		ID:        uuid.New(),
		Name:      data.Name,
		Email:     data.Email,
		Password:  data.Password, // In real app, this should be hashed
		Role:      data.Role,
		Provider:  data.Provider,
		Photo:     data.Photo,
		Verified:  data.Verified,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to repository
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	return nil
}

// UpdateUser updates an existing user
func (s *service) UpdateUser(id string, data *UpdateUserData) error {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user ID format")
	}

	// Business rule validations
	if data.Name == "" {
		return errors.New("name is required")
	}

	if data.Email == "" {
		return errors.New("email is required")
	}

	// Check if user exists
	existingUser, err := s.repo.GetUserByID(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if email is already taken by another user
	if data.Email != existingUser.Email {
		emailUser, err := s.repo.GetUserByEmail(data.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if emailUser != nil && emailUser.ID != existingUser.ID {
			return errors.New("email is already taken by another user")
		}
	}

	// Set default values if empty
	if data.Role == "" {
		data.Role = "user"
	}

	if data.Photo == "" {
		data.Photo = "default.png"
	}

	// Update user entity
	updatedUser := &User{
		Name:      data.Name,
		Role:      data.Role,
		Photo:     data.Photo,
		UpdatedAt: time.Now(),
	}

	// Save to repository
	if err := s.repo.UpdateUser(id, updatedUser); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *service) DeleteUser(id string) error {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user ID format")
	}

	// Check if user exists
	_, err := s.repo.GetUserByID(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Soft delete user
	if err := s.repo.DeleteUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

// RestoreUser restores a soft deleted user
func (s *service) RestoreUser(id string) (*UserResponse, error) {
	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Check if user exists (including soft deleted)
	user, err := s.repo.GetUserByID(id, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if user is actually deleted
	if user.DeletedAt == nil {
		return nil, errors.New("user is not deleted")
	}

	// Restore user
	if err := s.repo.RestoreUser(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Return restored user
	restoredUser, err := s.repo.GetUserByID(id, false)
	if err != nil {
		return nil, err
	}

	return restoredUser, nil
}

// CalculatePagination calculates pagination metadata
func (s *service) CalculatePagination(total int64, page, perPage int) int {
	if perPage <= 0 {
		perPage = 10
	}
	return int(math.Ceil(float64(total) / float64(perPage)))
}
