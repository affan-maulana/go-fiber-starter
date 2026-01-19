package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user domain entity (pure business model)
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      string
	Provider  string
	Photo     string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// CreateUserData represents user creation data for domain layer
type CreateUserData struct {
	Name     string
	Email    string
	Password string
	Role     string
	Provider string
	Photo    string
	Verified bool
}

// UpdateUserData represents user update data for domain layer
type UpdateUserData struct {
	Name     string
	Email    string
	Role     string
	Photo    string
	Verified bool
}

// UserModel represents the database model with GORM tags (infrastructure concern)
type UserModel struct {
	ID        *uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name      string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string     `gorm:"type:varchar(100);not null"`
	Role      string     `gorm:"type:varchar(50);default:'user';not null"`
	Provider  string     `gorm:"type:varchar(50);default:'local';not null"`
	Photo     string     `gorm:"type:text;default:'default.png';not null"`
	Verified  bool       `gorm:"not null;default:false"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	DeletedAt *time.Time `gorm:"index"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}
