package auth

import (
	"errors"
	"testing"

	"github.com/golang-fiber-jwt/internal/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) CreateUser(user *user.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByID(id string) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

// Test SignUp Service - Success
func TestService_SignUp_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	signUpData := &SignUpData{
		Name:            "John Doe",
		Email:           "john@example.com",
		Password:        "password123",
		PasswordConfirm: "password123",
		Photo:           "photo.jpg",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*user.User")).Return(nil)

	user, err := service.SignUp(signUpData)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, "user", user.Role)
	assert.Equal(t, "local", user.Provider)
	assert.NotEqual(t, "password123", user.Password) // Should be hashed
	mockRepo.AssertExpectations(t)
}

// Test SignUp Service - Password Mismatch
func TestService_SignUp_PasswordMismatch(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	signUpData := &SignUpData{
		Name:            "John Doe",
		Email:           "john@example.com",
		Password:        "password123",
		PasswordConfirm: "different123",
		Photo:           "photo.jpg",
	}

	user, err := service.SignUp(signUpData)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "passwords do not match", err.Error())
}

// Test SignUp Service - Validation Errors
func TestService_SignUp_ValidationErrors(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	tests := []struct {
		name          string
		signUpData    *SignUpData
		expectedError string
	}{
		{
			name: "Empty Name",
			signUpData: &SignUpData{
				Name:            "",
				Email:           "john@example.com",
				Password:        "password123",
				PasswordConfirm: "password123",
			},
			expectedError: "name is required",
		},
		{
			name: "Empty Email",
			signUpData: &SignUpData{
				Name:            "John Doe",
				Email:           "",
				Password:        "password123",
				PasswordConfirm: "password123",
			},
			expectedError: "email is required",
		},
		{
			name: "Empty Password",
			signUpData: &SignUpData{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "",
				PasswordConfirm: "",
			},
			expectedError: "password is required",
		},
		{
			name: "Password Too Short",
			signUpData: &SignUpData{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "short",
				PasswordConfirm: "short",
			},
			expectedError: "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.SignUp(tt.signUpData)
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Equal(t, tt.expectedError, err.Error())
		})
	}
}

// Test SignUp Service - Duplicate Email
func TestService_SignUp_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	signUpData := &SignUpData{
		Name:            "John Doe",
		Email:           "john@example.com",
		Password:        "password123",
		PasswordConfirm: "password123",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*auth.User")).
		Return(errors.New("duplicate key value violates unique constraint"))

	user, err := service.SignUp(signUpData)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user with that email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test SignIn Service - Success
func TestService_SignIn_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	// Create a user with hashed password
	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy" // "password123"
	existingUser := &user.User{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("GetUserByEmail", "john@example.com").Return(existingUser, nil)

	token, user, err := service.SignIn("john@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

// Test SignIn Service - User Not Found
func TestService_SignIn_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	mockRepo.On("GetUserByEmail", "notfound@example.com").Return(nil, errors.New("record not found"))

	token, user, err := service.SignIn("notfound@example.com", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test SignIn Service - Invalid Password
func TestService_SignIn_InvalidPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy" // "password123"
	existingUser := &user.User{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: hashedPassword,
	}

	mockRepo.On("GetUserByEmail", "john@example.com").Return(existingUser, nil)

	token, user, err := service.SignIn("john@example.com", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
	assert.Equal(t, "invalid email or password", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test GetUserByID Service - Success
func TestService_GetUserByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	userID := uuid.New().String()
	expectedUser := &user.User{
		ID:    uuid.MustParse(userID),
		Name:  "John Doe",
		Email: "john@example.com",
		Role:  "user",
	}

	mockRepo.On("GetUserByID", userID).Return(expectedUser, nil)

	user, err := service.GetUserByID(userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

// Test GetUserByID Service - User Not Found
func TestService_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewAuthService(mockRepo)

	userID := uuid.New().String()

	mockRepo.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	user, err := service.GetUserByID(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}
