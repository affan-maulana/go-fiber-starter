# Clean Architecture Implementation

## Project Structure

```
.
├── cmd/
│   └── main.go                          # Application entry point & DI wiring
│
├── config/
│   ├── db.go                            # Database configuration
│   ├── loadEnv.go                       # Environment loader
│   └── migrate.go                       # Database migrations
│
├── internal/
│   ├── auth/                            # Auth domain module
│   │   ├── auth_entity.go               # Domain entities + UserModel (GORM tags)
│   │   ├── auth_dto.go                  # HTTP DTOs (JSON/validation tags)
│   │   ├── auth_repository.go           # Repository interface + implementation
│   │   ├── auth_service.go              # Pure business logic (no framework deps)
│   │   └── auth_handler.go              # HTTP handlers (maps DTOs ↔ entities)
│   │
│   └── middleware/
│       └── deserialize-user.go          # JWT authentication middleware
│
├── pkg/                                 # Shared utilities
│   ├── validator/
│   │   └── validator.go                 # Reusable validation logic
│   ├── response/
│   │   └── response.go                  # Response formatters
│   └── hashing/
│       └── hashing.go                   # Hashing utilities
│
├── routes/
│   ├── routes.go                        # Main router setup
│   ├── auth.routes.go                   # Auth routes configuration
│   └── user.routes.go                   # User routes configuration
│
├── migrations/                          # Database migration files
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
│
├── docs/                                # Documentation
│   └── clean-architecture.md            # This file
│
├── docker-compose.yml
├── .env                                 # Environment variables
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Architecture Overview

This project follows **strict package-by-feature clean architecture** principles:

### Directory Structure

```
/internal/
  /auth/                    # Auth domain package
    model.go               # Pure domain entities (no framework tags)
    repository.go          # Repository interface only
    service.go             # Pure business logic (no framework imports)
    handler.go             # HTTP handlers (Fiber allowed)
    
  /infra/                  # Infrastructure implementations
    /postgresql/
      auth_repository.go   # GORM implementation of auth.Repository

/transport/                # Transport layer (HTTP concerns)
  /http/
    /auth/
      request.go           # HTTP request DTOs with JSON tags
      response.go          # HTTP response DTOs with JSON tags

/pkg/                      # Shared utilities
  /validator/
    validator.go           # Reusable validation logic
```

## Layer Responsibilities

### 1. Domain Layer (`/internal/<domain>/`)

Pure business logic with zero framework dependencies.

## Architecture Rules (STRICTLY ENFORCED)

### 1. Domain Packages (`/internal/<domain>/`)

**Each domain contains exactly 5 files:**
- `<domain>_entity.go` - Domain entities + Database models (GORM tags allowed for DB models)
- `<domain>_dto.go` - HTTP request/response DTOs (JSON/validation tags allowed)
- `<domain>_repository.go` - Repository interface + GORM implementation
- `<domain>_service.go` - Business logic (no framework imports)
- `<domain>_handler.go` - HTTP handlers (Fiber allowed, maps DTOs ↔ entities)

**Domain module organization:**
- **Entity file**: Domain entities (pure) + Database models (GORM tags allowed)
- **DTO file**: HTTP DTOs with JSON/validation tags
- **Repository file**: Interface + GORM implementation
- **Service file**: Pure business logic (NO framework imports)
- **Handler file**: HTTP layer (Fiber allowed, maps DTOs ↔ entities)

**Service layer MUST be pure:**
- ✅ Go stdlib only
- ✅ Domain entities only
- ❌ NO JSON tags
- ❌ NO Fiber imports
- ❌ NO GORM imports
- ❌ NO framework dependencies

### 2. Shared Utilities (`/pkg/`)

**Reusable utilities across all domains:**
- `/pkg/validator/` - Validation logic (used by all DTO layers)
- `/pkg/response/` - Response formatters
- `/pkg/hashing/` - Hashing utilities
- ✅ Domain-agnostic
- ✅ No business logic
- ✅ Stateless utilities only

### 5. Dependency Direction

```
Transport → Handler → Service → Repository (interface)
   ↓           ↓         ↓            ↑
   └───────────┴─────────┴────────────┘
                      |
              Infra implements
```

**Critical Rules:**
- Domain NEVER imports infra
- Domain NEVER imports transport
- Transport can import domain (for handler usage)
- Infra implements domain interfaces
- Handler maps between transport DTOs and domain models

### 4. Dependency Injection (main.go)

All wiring happens in `main.go`:

```go
// 1. Instantiate repository (now in same module)
db := config.SetupDatabase()

// 2. Inject into service
authService := auth.NewService(db) // Service creates its own repository

// 3. Inject into handler
authHandler := auth.NewHandler(authService)

// 4. Pass to routes
routes.SetupRoutes(app, authHandler)
```

## Implementation Details

### Auth Domain Flow

**1. Model (model.go)** - Pure domain entities
```go
type User struct {
    ID        uuid.UUID
    Name      string
    Email     string
    Password  string
    // ... no tags!
}

type SignUpData struct {
    Name            string
    Email           string
    Password        string
    PasswordConfirm string
    Photo           string
}
```

**2. Repository Interface (repository.go)** - Contract only
```go
type Repository interface {
    GetUserByEmail(email string) (*User, error)
    CreateUser(user *User) error
}
```

**3. Service (service.go)** - Pure business logic
```go
type Service interface {
    SignUp(data *SignUpData) (*User, error)
    SignIn(email, password string) (string, *User, error)
}

func (s *service) SignUp(data *SignUpData) (*User, error) {
    // Validate business rules
    if data.Password != data.PasswordConfirm {
        return nil, errors.New("passwords do not match")
    }
    // Hash password
    // Create user
    // No HTTP, no DB framework code in service
    return s.repo.CreateUser(user)
}
```

**5. Handler (auth_handler.go)** - HTTP layer with DTO mapping
```go
func (h *Handler) SignUpUser(c *fiber.Ctx) error {
    // Parse request DTO
    var req SignUpRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request")
    }
    
    // Map DTO → domain
    signUpData := &SignUpData{
        Name:            req.Name,
        Email:           req.Email,
        Password:        req.Password,
        PasswordConfirm: req.PasswordConfirm,
        Photo:           req.Photo,
    }
    
    // Call service
    user, err := h.service.SignUp(signUpData)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }
    
    // Map domain → response DTO
    userResponse := UserResponse{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        // ...
    }
    
    return response.Created(c, UserDataResponse{User: userResponse})
}
    
    // Map Domain → Response DTO
    response := h.mapUserToResponse(user)
    
    return c.JSON(response)
}
```

### Benefits

1. **Testability**: Each domain can be tested independently
2. **No Framework Lock-in**: Business logic has zero framework dependencies
3. **Clear Boundaries**: Domain vs Infrastructure separation within each module
4. **Single Responsibility**: Each layer has one job
5. **Dependency Inversion**: High-level policies don't depend on low-level details
6. **Modular**: Each domain module is self-contained

## Testing Strategy

### Test Structure

```
/internal/auth/
  service_test.go         # Unit tests for business logic
  
/transport/http/auth/
  handler_integration_test.go  # Integration tests for HTTP layer
```

### 1. Unit Tests (Domain Layer)

**Location:** `/internal/<domain>/service_test.go`

**Purpose:** Test pure business logic with mocked dependencies

```go
package auth

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockRepository mocks the Repository interface
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetUserByEmail(email string) (*User, error) {
    args := m.Called(email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// Test business logic without dependencies
func TestService_SignUp_Success(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    
    signUpData := &SignUpData{
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: "password123",
        PasswordConfirm: "password123",
    }
    
    mockRepo.On("CreateUser", mock.AnythingOfType("*auth.User")).Return(nil)
    
    user, err := service.SignUp(signUpData)
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
    mockRepo.AssertExpectations(t)
}
```

**What to test:**
- ✅ Business logic validation
- ✅ Error handling
- ✅ Edge cases
- ✅ Password hashing
- ✅ Data transformations
- ❌ No HTTP concerns
- ❌ No database queries

### 2. Integration Tests (Transport Layer)

**Location:** `/transport/http/<domain>/handler_integration_test.go`

**Purpose:** Test HTTP handlers with mocked services

**Important:** Use `package <domain>_test` to avoid import cycles

```go
package auth_test

import (
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    authDomain "github.com/golang-fiber-jwt/internal/auth"
    authTransport "github.com/golang-fiber-jwt/transport/http/auth"
)

// MockService mocks the Service interface
type MockService struct {
    mock.Mock
}

func (m *MockService) SignUp(data *authDomain.SignUpData) (*authDomain.User, error) {
    args := m.Called(data)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*authDomain.User), args.Error(1)
}

// Test HTTP layer
func TestSignUpUser_Success(t *testing.T) {
    app := fiber.New()
    mockService := new(MockService)
    handler := authDomain.NewHandler(mockService)
    
    app.Post("/register", handler.SignUpUser)
    
    requestBody := authTransport.SignUpRequest{
        Name:     "John Doe",
        Email:    "john@example.com",
        Password: "password123",
        PasswordConfirm: "password123",
    }
    
    expectedUser := &authDomain.User{...}
    mockService.On("SignUp", mock.AnythingOfType("*auth.SignUpData")).Return(expectedUser, nil)
    
    // Make HTTP request and assert response
    resp, _ := app.Test(req)
    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
```

**What to test:**
- ✅ Request parsing
- ✅ Validation errors
- ✅ Response formatting
- ✅ HTTP status codes
- ✅ DTO mapping
- ❌ No business logic (use mocks)

### 3. Running Tests

```bash
# Run all tests
go test ./...

# Run domain unit tests
go test ./internal/auth/...

# Run integration tests
go test ./transport/http/auth/...

# Run with coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### 4. Test Organization

**Unit Tests (Domain):**
- Test each service method
- Mock all external dependencies
- Focus on business logic correctness
- Fast execution (no I/O)

**Integration Tests (Transport):**
- Test HTTP request/response flow
- Mock service layer
- Verify DTO transformations
- Test middleware integration

**Example Test Coverage:**
```
internal/auth/service_test.go:
  ✓ TestService_SignUp_Success
  ✓ TestService_SignUp_PasswordMismatch
  ✓ TestService_SignUp_ValidationErrors
  ✓ TestService_SignUp_DuplicateEmail
  ✓ TestService_SignIn_Success
  ✓ TestService_SignIn_UserNotFound
  ✓ TestService_SignIn_InvalidPassword
  ✓ TestService_GetUserByID_Success
  ✓ TestService_GetUserByID_NotFound

transport/http/auth/handler_integration_test.go:
  ✓ TestSignUpUser_Success
  ✓ TestSignUpUser_ValidationError
  ✓ TestSignUpUser_DuplicateEmail
  ✓ TestSignInUser_Success
  ✓ TestSignInUser_InvalidCredentials
  ✓ TestLogoutUser_Success
  ✓ TestGetMe_Success
  ✓ TestGetMe_Unauthorized
  ✓ TestGetMe_UserNotFound
```

## Adding New Features (Step-by-Step Guide)

### Complete Workflow: Adding a New Module

Follow these steps to add a new domain module (e.g., `product`, `order`, `customer`) while maintaining clean architecture:

---

### Example: Add "Product" Domain

#### Step 1: Create Domain Package

**Location:** `/internal/product/`

**1. Create domain package** (`/internal/product/`)

```go
// model.go - Pure domain
type Product struct {
    ID          uuid.UUID
    Name        string
    Price       float64
    Stock       int
    CreatedAt   time.Time
}

// repository.go - Interface only
type Repository interface {
    CreateProduct(product *Product) error
    GetProductByID(id string) (*Product, error)
}

// service.go - Business logic
type Service interface {
    CreateProduct(data *CreateProductData) (*Product, error)
}

// handler.go - HTTP layer
type Handler struct {
    service Service
}
```

**2. Create transport DTOs** (`/transport/http/product/`)

```go
// request.go
type CreateProductRequest struct {
    Name  string  `json:"name" validate:"required"`
    Price float64 `json:"price" validate:"required,gt=0"`
    Stock int     `json:"stock" validate:"required,gte=0"`
}

// response.go
type ProductResponse struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    Stock     int       `json:"stock"`
    CreatedAt time.Time `json:"created_at"`
}
```

**3. Create infrastructure** (`/internal/infra/postgresql/`)

```go
// product_repository.go
type ProductModel struct {
    ID        *uuid.UUID `gorm:"type:uuid;primary_key"`
    Name      string     `gorm:"type:varchar(255);not null"`
    Price     float64    `gorm:"type:decimal(10,2);not null"`
    Stock     int        `gorm:"type:int;not null"`
    CreatedAt *time.Time `gorm:"not null;default:now()"`
}

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.Repository {
    return &productRepository{db: db}
}
```

**4. Wire in main.go**

```go
// Instantiate infrastructure
productRepo := postgresql.NewProductRepository(config.DB)

// Inject into domain
productService := product.NewService(productRepo)
productHandler := product.NewHandler(productService)

// Pass to routes
routes.SetupRoutes(app, authHandler, productHandler)
```

**5. Create routes** (`/routes/product.routes.go`)

```go
func ProductRoutes(router fiber.Router, handler *product.Handler) {
    router.Post("/products", handler.CreateProduct)
    router.Get("/products/:id", handler.GetProduct)
}
```

**No cross-domain imports! Each domain is independent.**

---

### Quick Reference: New Module Checklist

When creating a new module, follow this checklist:

#### 📁 Step 1: Domain Layer (`/internal/<domain>/`)

```bash
# Create domain directory
mkdir -p internal/<domain>
```

**Create these 4 files:**

1️⃣ **`model.go`** - Pure domain entities
```go
package <domain>

import "github.com/google/uuid"

// Domain entity - NO JSON/GORM tags
type <Entity> struct {
    ID        uuid.UUID
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Domain data transfer objects (internal use)
type Create<Entity>Data struct {
    Name string
}
```

2️⃣ **`repository.go`** - Interface ONLY
```go
package <domain>

// Repository interface - no implementation here
type Repository interface {
    Create<Entity>(entity *<Entity>) error
    Get<Entity>ByID(id string) (*<Entity>, error)
    Update<Entity>(entity *<Entity>) error
    Delete<Entity>(id string) error
}
```

3️⃣ **`service.go`** - Pure business logic
```go
package <domain>

import "errors"

type Service interface {
    Create<Entity>(data *Create<Entity>Data) (*<Entity>, error)
    Get<Entity>ByID(id string) (*<Entity>, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) Create<Entity>(data *Create<Entity>Data) (*<Entity>, error) {
    // Validate business rules
    if data.Name == "" {
        return nil, errors.New("name is required")
    }
    
    // Create entity
    entity := &<Entity>{
        ID:        uuid.New(),
        Name:      data.Name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Persist
    if err := s.repo.Create<Entity>(entity); err != nil {
        return nil, err
    }
    
    return entity, nil
}
```

4️⃣ **`handler.go`** - HTTP handlers with DTO mapping
```go
package <domain>

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-fiber-jwt/pkg/validator"
    transport "github.com/golang-fiber-jwt/transport/http/<domain>"
)

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create<Entity>(c *fiber.Ctx) error {
    var req transport.Create<Entity>Request
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "fail",
            "message": err.Error(),
        })
    }
    
    // Validate
    if errors := validator.ValidateStruct(req); errors != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "fail",
            "errors": errors,
        })
    }
    
    // Map DTO → Domain
    data := &Create<Entity>Data{
        Name: req.Name,
    }
    
    // Call service
    entity, err := h.service.Create<Entity>(data)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status": "fail",
            "message": err.Error(),
        })
    }
    
    // Map Domain → Response DTO
    response := transport.<Entity>Response{
        ID:        entity.ID,
        Name:      entity.Name,
        CreatedAt: entity.CreatedAt,
    }
    
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status": "success",
        "data": fiber.Map{"<entity>": response},
    })
}
```

#### 📁 Step 2: Transport Layer (`/transport/http/<domain>/`)

```bash
# Create transport directory
mkdir -p transport/http/<domain>
```

**Create these 2 files:**

1️⃣ **`request.go`** - HTTP request DTOs with validation
```go
package <domain>

// Request DTO - JSON and validation tags allowed
type Create<Entity>Request struct {
    Name string `json:"name" validate:"required"`
}

type Update<Entity>Request struct {
    Name string `json:"name" validate:"required"`
}
```

2️⃣ **`response.go`** - HTTP response DTOs
```go
package <domain>

import (
    "time"
    "github.com/google/uuid"
)

// Response DTO - JSON tags allowed
type <Entity>Response struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Errors  interface{} `json:"errors,omitempty"`
}
```

#### 📁 Step 3: Infrastructure Layer (`/internal/infra/postgresql/`)

**Create:** `<domain>_repository.go`

```go
package postgresql

import (
    "gorm.io/gorm"
    "github.com/google/uuid"
    "<yourmodule>/internal/<domain>"
)

// Database model with GORM tags
type <Entity>Model struct {
    ID        *uuid.UUID `gorm:"type:uuid;primary_key"`
    Name      string     `gorm:"type:varchar(255);not null"`
    CreatedAt *time.Time `gorm:"not null;default:now()"`
    UpdatedAt *time.Time `gorm:"not null;default:now()"`
}

func (m *<Entity>Model) TableName() string {
    return "<entities>"
}

// Repository implementation
type <domain>Repository struct {
    db *gorm.DB
}

func New<Domain>Repository(db *gorm.DB) <domain>.Repository {
    return &<domain>Repository{db: db}
}

func (r *<domain>Repository) Create<Entity>(entity *<domain>.<Entity>) error {
    model := toModel(entity)
    if err := r.db.Create(&model).Error; err != nil {
        return err
    }
    return nil
}

func (r *<domain>Repository) Get<Entity>ByID(id string) (*<domain>.<Entity>, error) {
    var model <Entity>Model
    if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
        return nil, err
    }
    return toDomain(&model), nil
}

// Mapper: Domain → Database Model
func toModel(entity *<domain>.<Entity>) *<Entity>Model {
    return &<Entity>Model{
        ID:        &entity.ID,
        Name:      entity.Name,
        CreatedAt: &entity.CreatedAt,
        UpdatedAt: &entity.UpdatedAt,
    }
}

// Mapper: Database Model → Domain
func toDomain(model *<Entity>Model) *<domain>.<Entity> {
    return &<domain>.<Entity>{
        ID:        *model.ID,
        Name:      model.Name,
        CreatedAt: *model.CreatedAt,
        UpdatedAt: *model.UpdatedAt,
    }
}
```

#### 📁 Step 4: Wire Dependencies (`/internal/container/container.go`)

```go
// Update container.go
package container

import (
    "gorm.io/gorm"
    "<yourmodule>/internal/auth"
    "<yourmodule>/internal/<domain>"
    "<yourmodule>/internal/infra/postgresql"
)

type Container struct {
    AuthHandler    *auth.Handler
    <Domain>Handler *<domain>.Handler
}

func NewContainer(db *gorm.DB) *Container {
    // Auth (existing)
    authRepo := postgresql.NewAuthRepository(db)
    authService := auth.NewService(authRepo)
    authHandler := auth.NewHandler(authService)
    
    // New domain
    <domain>Repo := postgresql.New<Domain>Repository(db)
    <domain>Service := <domain>.NewService(<domain>Repo)
    <domain>Handler := <domain>.NewHandler(<domain>Service)
    
    return &Container{
        AuthHandler:    authHandler,
        <Domain>Handler: <domain>Handler,
    }
}
```

#### 📁 Step 5: Create Routes (`/routes/<domain>.routes.go`)

```go
package routes

import (
    "github.com/gofiber/fiber/v2"
    "<yourmodule>/internal/<domain>"
)

func <Domain>Routes(router fiber.Router, handler *<domain>.Handler) {
    <domain> := router.Group("/<entities>")
    
    <domain>.Post("/", handler.Create<Entity>)
    <domain>.Get("/:id", handler.Get<Entity>ByID)
    <domain>.Put("/:id", handler.Update<Entity>)
    <domain>.Delete("/:id", handler.Delete<Entity>)
}
```

#### 📁 Step 6: Register Routes (`/routes/routes.go`)

```go
func SetupRoutes(app *fiber.App, container *container.Container) {
    api := app.Group("/api")
    
    // Existing routes
    AuthRoutes(api, container.AuthHandler)
    
    // New domain routes
    <Domain>Routes(api, container.<Domain>Handler)
}
```

#### 📁 Step 7: Create Tests

**Unit Tests:** `/internal/<domain>/service_test.go`

```go
package <domain>

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create<Entity>(entity *<Entity>) error {
    args := m.Called(entity)
    return args.Error(0)
}

func TestService_Create<Entity>_Success(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    
    data := &Create<Entity>Data{Name: "Test"}
    mockRepo.On("Create<Entity>", mock.AnythingOfType("*<domain>.<Entity>")).Return(nil)
    
    entity, err := service.Create<Entity>(data)
    
    assert.NoError(t, err)
    assert.NotNil(t, entity)
    mockRepo.AssertExpectations(t)
}
```

**Integration Tests:** `/transport/http/<domain>/handler_integration_test.go`

```go
package <domain>_test

import (
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    domain<Domain> "<yourmodule>/internal/<domain>"
    transport<Domain> "<yourmodule>/transport/http/<domain>"
)

type MockService struct {
    mock.Mock
}

func TestCreate<Entity>_Success(t *testing.T) {
    app := fiber.New()
    mockService := new(MockService)
    handler := domain<Domain>.NewHandler(mockService)
    
    // Test HTTP endpoint
}
```

#### 📁 Step 8: Create Migrations (if needed)

```bash
# Create migration files
touch migrations/000002_create_<entities>_table.up.sql
touch migrations/000002_create_<entities>_table.down.sql
```

**Up migration:**
```sql
CREATE TABLE IF NOT EXISTS <entities> (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Down migration:**
```sql
DROP TABLE IF EXISTS <entities>;
```

#### ✅ Final Verification Checklist

Before committing, verify:

```bash
# 1. Build succeeds
go build ./cmd/main.go

# 2. No JSON tags in domain models
grep -r "json:" internal/<domain>/model.go  # Should be empty

# 3. No framework imports in service
grep -r "gorm\|fiber" internal/<domain>/service.go internal/<domain>/model.go  # Should be empty

# 4. Tests pass
go test ./internal/<domain>/...
go test ./transport/http/<domain>/...

# 5. No circular imports
go build -v ./...
```

**Architecture compliance:**
- [ ] Domain models have NO tags
- [ ] Service has NO framework imports
- [ ] Repository is interface-only in domain
- [ ] DTOs are in transport layer
- [ ] Infrastructure implements domain interfaces
- [ ] Handler maps DTOs ↔ Domain models
- [ ] Tests follow the pattern (unit + integration)
- [ ] Container wires all dependencies

---

### Summary

Each new module requires:
1. **Domain** (4 files): model, repository, service, handler
2. **Transport** (2 files): request, response
3. **Infrastructure** (1 file): repository implementation
4. **Routes** (1 file): route definitions
5. **Container** (update): wire dependencies
6. **Tests** (2 files): unit tests + integration tests
7. **Migrations** (optional): database schema

**Total:** ~10 files per domain, all following clean architecture principles.

## Migration Notes

### What Changed from Old Structure

**Before (Modules Architecture):**
```
/internal/modules/
  /entity/
    user.go              # Global entities
  /auth/
    auth.handler.go      # Mixed concerns
    auth.service.go      # Had JWT/config imports
    auth.repository.go   # Interface + GORM in same file
    /dto/
      user_request.go    # DTOs inside domain
      user_response.go
```

**After (Clean Architecture):**
```
/internal/
  /auth/
    model.go            # Pure domain (no tags)
    repository.go       # Interface only
    service.go          # Pure business logic
    handler.go          # HTTP layer
  /infra/postgresql/
    auth_repository.go  # GORM implementation

/transport/http/auth/
  request.go            # HTTP DTOs with JSON tags
  response.go

/pkg/validator/
  validator.go          # Shared validation
```

### Key Improvements

1. ✅ **Separation of Concerns**
   - Domain models have NO framework tags
   - DTOs separated from domain
   - Infrastructure separated from business logic

2. ✅ **DTO Migration**
   - Moved from `/internal/auth/dto/` to `/transport/http/auth/`
   - Clear separation: transport vs domain
   - JSON tags only in transport layer

3. ✅ **No Framework Lock-in**
   - Domain has zero framework dependencies
   - Can swap Fiber → Echo without touching domain
   - Can swap GORM → sqlx without touching domain

4. ✅ **Testability**
   - Mock interfaces easily
   - Test business logic without DB
   - Test without HTTP framework

5. ✅ **Reusability**
   - Validator shared across all domains
   - No code duplication
   - Consistent patterns

## Validation & Best Practices

### Verify Clean Boundaries

Run these commands to ensure architecture compliance:

```bash
# Build should succeed
go build ./cmd/main.go

# Check domain has no JSON tags
grep -r "json:" internal/auth/model.go  # Should be empty

# Check domain has no framework imports
grep -r "gorm\|fiber" internal/auth/service.go internal/auth/model.go  # Should be empty

# Check infra implements domain
grep "auth.Repository" internal/infra/postgresql/auth_repository.go  # Should exist

# Check no circular imports
go build -v ./...
```

### Architecture Checklist

When adding new code, verify:

- [ ] Domain models have NO tags (no `json:`, no `gorm:`)
- [ ] Domain service has NO framework imports (no Fiber, GORM, Redis)
- [ ] Repository is interface-only in domain
- [ ] Infrastructure implements domain interfaces
- [ ] DTOs are in `/transport/http/<domain>/`
- [ ] JSON tags only in transport layer
- [ ] GORM tags only in infrastructure layer
- [ ] Handler maps between DTOs and domain models
- [ ] No cross-domain imports
- [ ] Validation uses shared `/pkg/validator/`

### Common Violations to Avoid

❌ **Adding JSON tags to domain:**
```go
// WRONG - in internal/auth/model.go
type User struct {
    ID   uuid.UUID `json:"id"`  // ❌ NO JSON TAGS IN DOMAIN
    Name string    `json:"name"` // ❌
}
```

❌ **Importing framework in service:**
```go
// WRONG - in internal/auth/service.go
import "github.com/gofiber/fiber/v2"  // ❌ NO FIBER IN SERVICE
```

❌ **DTOs in domain package:**
```go
// WRONG - creating internal/auth/dto/
// DTOs must be in /transport/http/auth/  // ❌
```

❌ **GORM in domain:**
```go
// WRONG - in internal/auth/repository.go
type authRepository struct {
    db *gorm.DB  // ❌ NO GORM IN DOMAIN
}
```

✅ **Correct patterns shown in examples above**

## Future Domains

When adding new domains (e.g., `product`, `order`, `customer`):

### Step-by-Step Process

1. **Create domain package** `/internal/<domain>/`
   - `model.go` - Pure entities
   - `repository.go` - Interface only
   - `service.go` - Business logic
   - `handler.go` - HTTP handlers

2. **Create transport DTOs** `/transport/http/<domain>/`
   - `request.go` - Request DTOs with validation tags
   - `response.go` - Response DTOs with JSON tags

3. **Create infrastructure** `/internal/infra/postgresql/`
   - `<domain>_repository.go` - GORM implementation

4. **Wire in `main.go`**
   ```go
   domainRepo := postgresql.NewDomainRepository(config.DB)
   domainService := domain.NewService(domainRepo)
   domainHandler := domain.NewHandler(domainService)
   ```

5. **Create routes** `/routes/<domain>.routes.go`

6. **Add to main router** in `routes/routes.go`

### Domain Independence

Each domain is a **vertical slice:**
- Independently testable
- Independently deployable (microservices future)
- No cross-domain imports
- Self-contained business logic

**Example:** Product domain never imports Order domain. If they need to communicate, use:
- Events
- Shared interfaces
- API calls
- Message queues

This ensures loose coupling and high cohesion.
