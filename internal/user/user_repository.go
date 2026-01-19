package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AllowedSearchFields defines which fields can be searched
var AllowedSearchFields = map[string]bool{
	"name":  true,
	"email": true,
}

// Repository defines the interface for user data persistence
type Repository interface {
	// GetUsers retrieves users with filtering and pagination
	GetUsers(query ListUsersQuery) ([]UserResponse, int64, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(id string, includeDeleted bool) (*UserResponse, error)

	// GetUserByEmail retrieves a user by their email address
	GetUserByEmail(email string) (*UserResponse, error)

	// CreateUser creates a new user in the system
	CreateUser(user *User) error

	// UpdateUser updates an existing user
	UpdateUser(id string, user *User) error

	// DeleteUser soft deletes a user
	DeleteUser(id string) error

	// RestoreUser restores a soft deleted user
	RestoreUser(id string) error

	// HardDeleteUser permanently deletes a user
	HardDeleteUser(id string) error
}

// userRepository implements Repository interface with GORM
type userRepository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

// GetUsers retrieves users with filtering and pagination
func (r *userRepository) GetUsers(query ListUsersQuery) ([]UserResponse, int64, error) {
	var models []UserResponse
	var total int64

	// Set default pagination
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PerPage <= 0 {
		query.PerPage = 10
	}

	// Build base query
	db := r.db.Model(&UserModel{})

	// Include soft deleted records if requested
	if query.ShowDeleted {
		db = db.Unscoped()
	}

	// Apply dynamic search filter
	if query.Search != "" && query.SearchBy != "" {
		// Validate SearchBy field is allowed
		if AllowedSearchFields[query.SearchBy] {
			searchTerm := fmt.Sprintf("%%%s%%", query.Search)
			db = db.Where(fmt.Sprintf("%s ILIKE ?", query.SearchBy), searchTerm)
		}
	}

	// Apply standard filters
	if query.Role != "" {
		db = db.Where("role = ?", query.Role)
	}

	if query.Provider != "" {
		db = db.Where("provider = ?", query.Provider)
	}

	if query.Verified != nil {
		db = db.Where("verified = ?", *query.Verified)
	}

	// Count total records
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (query.Page - 1) * query.PerPage
	err := db.Select("id, name, email, role, photo, created_at").
		Offset(offset).
		Limit(query.PerPage).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	// Convert to domain models
	users := make([]UserResponse, len(models))
	for i, model := range models {
		users[i] = *toDomain(&model)
	}

	return users, total, nil
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(id string, includeDeleted bool) (*UserResponse, error) {
	var model UserResponse

	db := r.db
	if includeDeleted {
		db = db.Unscoped()
	}

	result := db.Where("id = ?", id).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}

	return toDomain(&model), nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepository) GetUserByEmail(email string) (*UserResponse, error) {
	var model UserResponse
	result := r.db.Where("email = ?", email).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return toDomain(&model), nil
}

// CreateUser creates a new user
func (r *userRepository) CreateUser(user *User) error {
	model := toModel(user)
	result := r.db.Create(&model)
	if result.Error != nil {
		return result.Error
	}

	// Update user with generated values
	if model.ID != nil {
		user.ID = *model.ID
	}
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(id string, user *User) error {
	model := toModel(user)
	model.UpdatedAt = time.Now()

	result := r.db.Where("id = ?", id).Updates(&model)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	user.UpdatedAt = model.UpdatedAt
	return nil
}

// DeleteUser soft deletes a user
func (r *userRepository) DeleteUser(id string) error {
	result := r.db.Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// RestoreUser restores a soft deleted user
func (r *userRepository) RestoreUser(id string) error {
	result := r.db.Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// HardDeleteUser permanently deletes a user
func (r *userRepository) HardDeleteUser(id string) error {
	result := r.db.Unscoped().Where("id = ?", id).Delete(&UserModel{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// toDomain converts database model to domain model
func toDomain(model *UserResponse) *UserResponse {
	user := &UserResponse{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Role:      model.Role,
		Provider:  model.Provider,
		Photo:     model.Photo,
		Verified:  model.Verified,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	if model.DeletedAt != nil {
		user.DeletedAt = model.DeletedAt
	}

	return user
}

// toModel converts domain model to database model
func toModel(user *User) *UserModel {
	model := &UserModel{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		Provider:  user.Provider,
		Photo:     user.Photo,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	// Set ID if it exists
	if user.ID != uuid.Nil {
		id := user.ID
		model.ID = &id
	}

	return model
}
