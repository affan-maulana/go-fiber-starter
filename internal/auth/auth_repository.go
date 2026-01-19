package auth

import (
	"github.com/golang-fiber-jwt/internal/user"
	"gorm.io/gorm"
)

// Repository defines the interface for auth data persistence
// This is a pure interface with no implementation details
// Infrastructure layer will implement this interface
type Repository interface {
	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(email string) (*user.User, error)

	// CreateUser creates a new user in the system
	CreateUser(user *user.User) error

	// GetUserByID retrieves a user by their ID
	GetUserByID(id string) (*user.User, error)
}

// authRepository implements Repository interface
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new MySQL auth repository
func NewAuthRepository(db *gorm.DB) Repository {
	return &authRepository{db: db}
}

// GetUserByEmail retrieves a user by email
func (r *authRepository) GetUserByEmail(email string) (*user.User, error) {
	var model user.User
	result := r.db.Where("email = ?", email).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}

// CreateUser creates a new user
func (r *authRepository) CreateUser(user *user.User) error {
	result := r.db.Create(user)
	return result.Error
}

// GetUserByID retrieves a user by ID
func (r *authRepository) GetUserByID(id string) (*user.User, error) {
	var model user.User
	result := r.db.Where("id = ?", id).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}
