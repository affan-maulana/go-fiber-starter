# Go Fiber JWT - Clean Architecture

Backend API with Go Fiber following strict clean architecture principles.

## Architecture

This project implements **package-by-feature clean architecture** with:
- ✅ Pure domain models (no framework dependencies)
- ✅ Separated transport layer (DTOs with JSON tags)
- ✅ Infrastructure abstraction (GORM implementations)
- ✅ Dependency injection
- ✅ Testable business logic

📖 **See [Clean Architecture Documentation](./docs/clean-architecture.md) for detailed architecture guide**

## Project Structure

```
.
├── cmd/main.go                  # Entry point & DI wiring
├── internal/
│   ├── auth/                    # Auth domain module
│   │   ├── auth_entity.go       # Domain entities + UserModel
│   │   ├── auth_dto.go          # HTTP DTOs (request/response)
│   │   ├── auth_repository.go   # Repository interface + implementation
│   │   ├── auth_service.go      # Business logic service
│   │   └── auth_handler.go      # HTTP handlers
│   └── middleware/              # HTTP middlewares
├── pkg/                         # Shared utilities
│   ├── validator/               # Validation utilities
│   ├── response/                # Response formatters
│   └── hashing/                 # Hashing utilities
├── routes/                      # Route configurations
└── migrations/                  # Database migrations
```

## Getting Started

### Prerequisites
- Go 1.19+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd golang-fiber-jwt
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**

Copy `.env.example` to `.env` and configure:
```env
POSTGRES_HOST=127.0.0.1
POSTGRES_PORT=6500
POSTGRES_USER=admin
POSTGRES_PASSWORD=password123
POSTGRES_DB=golang-fiber-jwt

JWT_SECRET=your-secret-key
JWT_EXPIRED_IN=60m
JWT_MAXAGE=60
```

4. **Start database (Docker)**
```bash
docker-compose up -d
```

5. **Run migrations**

Migrations run automatically on startup, or manually:
```bash
migrate -path ./migrations -database "postgres://admin:password123@localhost:6500/golang-fiber-jwt?sslmode=disable" up
```

### Running the Application

**Development mode:**
```bash
go run ./cmd/main.go
```

**Build and run:**
```bash
go build -o app ./cmd/main.go
./app
```

**Using Makefile:**
```bash
make run      # Run application
make build    # Build binary
make test     # Run tests
```

Server runs on `http://localhost:3334`

## API Endpoints

### Auth

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/logout` - User logout (requires auth)

### User

- `GET /api/users/me` - Get current user (requires auth)

### Health

- `GET /api/healthchecker` - Health check endpoint

## Architecture Principles

This project follows strict clean architecture rules:

### Domain Module (`/internal/<domain>/`)
- **`<domain>_entity.go`** - Pure domain entities + database models (GORM tags allowed for UserModel)
- **`<domain>_dto.go`** - HTTP DTOs with JSON/validation tags
- **`<domain>_repository.go`** - Repository interface + implementation
- **`<domain>_service.go`** - Pure business logic (no framework dependencies)
- **`<domain>_handler.go`** - HTTP handlers (Fiber allowed)

### Layer Separation Within Module
- **Entity layer** - Domain entities (pure) + Database models (GORM tags)
- **DTO layer** - HTTP request/response with JSON/validation tags
- **Repository layer** - Data access interface + implementation
- **Service layer** - Business logic (framework-independent) + calculations (pagination, etc.)
- **Handler layer** - HTTP concerns (maps DTOs ↔ entities) + concurrency coordination

### Benefits
- ✅ **Testable** - Mock interfaces easily
- ✅ **Maintainable** - Clear separation of concerns
- ✅ **Flexible** - Swap implementations without changing business logic
- ✅ **Scalable** - Independent domain packages

📖 **[Read full architecture guide →](./docs/clean-architecture.md)**

## Database Migrations

### Creating Migrations

Migrations follow the naming convention: `{version}_{description}.{up|down}.sql`

**Example:**
```bash
# Create new migration files
touch migrations/000002_add_products_table.up.sql
touch migrations/000002_add_products_table.down.sql
```

**Up migration** (`000002_add_products_table.up.sql`):
```sql
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Down migration** (`000002_add_products_table.down.sql`):
```sql
DROP TABLE IF EXISTS products;
```

### Running Migrations

**Automatic** - Migrations run on application startup

**Manual:**
```bash
# Apply all pending migrations
migrate -path ./migrations -database $DATABASE_URL up

# Rollback one migration
migrate -path ./migrations -database $DATABASE_URL down 1

# Check version
migrate -path ./migrations -database $DATABASE_URL version
```

## Development

### Quick Start: Creating a New Module

#### Option 1: Manual Steps (Traditional Approach)

Follow these steps to add a new domain module while maintaining clean architecture:

1. **Create Domain Package** - `mkdir -p internal/product`
2. **Create Transport DTOs** - `mkdir -p transport/http/product`
3. **Create Infrastructure** - `internal/infra/postgresql/product_repository.go`
4. **Wire Dependencies** - Update `/internal/container/container.go`
5. **Create Routes** - Create `routes/product.routes.go`
6. **Add Tests** - Unit and integration tests
7. **Create Migrations** - Database migrations if needed

#### Option 2: AI-Generated CRUD Module (Recommended)

Use this template prompt to generate a complete CRUD module with all necessary files:

---

**🤖 AI PROMPT TEMPLATE - COPY & EDIT THIS:**

```
Generate a complete CRUD module for Go Fiber with Clean Architecture.

**Project Context:**
- Module Name from go.mod: github.com/golang-fiber-jwt
- Go Version: 1.24+
- Framework: Go Fiber + GORM
- Available Helper Functions: 
  - handler.ParseAndValidate(c *fiber.Ctx, req interface{}) error
  - response.OK(), response.Created(), response.BadRequest(), etc.

**Module Details:**
- Module Name: [CHANGE_THIS: product, category, order, etc.]
- Package Path: internal/[MODULE_NAME]/
- File Naming: [MODULE_NAME]_entity.go, [MODULE_NAME]_dto.go, etc.

**Database Table Structure:** (Copy from PostgreSQL)
```
Column Name	#	Data type	Identity	Collation	Not Null	Default	Comment
id	1	uuid	[NULL]	[NULL]	true	gen_random_uuid()	[NULL]
name	2	varchar(100)	[NULL]	default	true	[NULL]	[NULL]
description	3	text	[NULL]	default	false	[NULL]	[NULL]
price	4	decimal(10,2)	[NULL]	[NULL]	true	[NULL]	[NULL]
created_at	5	timestamptz	[NULL]	[NULL]	true	now()	[NULL]
updated_at	6	timestamptz	[NULL]	[NULL]	true	now()	[NULL]
deleted_at	7	timestamptz	[NULL]	[NULL]	false	[NULL]	[NULL]
```

**Requirements:**
- Framework: Go Fiber + GORM + Clean Architecture
- CRUD Operations: List, Detail, Create, Update, Delete (soft delete if deleted_at exists)
- File Structure:
  - `internal/[MODULE_NAME]/[MODULE_NAME]_entity.go` - Domain entities + Database model with GORM tags
  - `internal/[MODULE_NAME]/[MODULE_NAME]_dto.go` - HTTP request/response DTOs (JSON/validation tags)
  - `internal/[MODULE_NAME]/[MODULE_NAME]_repository.go` - Repository interface + GORM implementation
  - `internal/[MODULE_NAME]/[MODULE_NAME]_service.go` - Business logic service (no framework dependencies)
  - `internal/[MODULE_NAME]/[MODULE_NAME]_handler.go` - HTTP handlers with manual DTO mapping
- Import Statements: Use "github.com/golang-fiber-jwt/pkg/response" and "github.com/golang-fiber-jwt/pkg/handler"
- Response: Use response.OK(), response.Created(), response.BadRequest(), response.NotFound(), response.InternalError()
- Validation: Add `validate` tags for required fields (password/email always required, others based on business needs)
- Soft Delete: Implement if `deleted_at` column exists using GORM soft delete with DeletedAt *time.Time field
- UUID Primary Key: Use `github.com/google/uuid` with proper nil checks
- Error Handling: Proper HTTP status codes with handleServiceError() method
- Handler Methods:
  - `List[MODULE_NAME]s(c *fiber.Ctx) error` (GET with pagination)
  - `Get[MODULE_NAME]ByID(c *fiber.Ctx) error` (GET by ID)
  - `Create[MODULE_NAME](c *fiber.Ctx) error` (POST)
  - `Update[MODULE_NAME](c *fiber.Ctx) error` (PUT)
  - `Delete[MODULE_NAME](c *fiber.Ctx) error` (DELETE)

**Handler Implementation Rules:**
- Use handler.ParseAndValidate(c, &req) for request parsing
- Manual DTO mapping (no generic helpers)
- Query parameters: Parse manually using c.Query() and strconv
- URL parameters: Use c.Params("id")
- Include strconv and sync imports for parameter parsing and concurrency
- Implement goroutines for parallel processing where beneficial (e.g., data fetching + pagination)
- Use channels and sync.WaitGroup for concurrent operations
- Business logic calculations (like pagination) should be in service layer

**Code Style:**
- Follow existing project conventions
- All files in single module package
- Repository interface + implementation in same file
- Domain entities separate from database models
- Manual mapping between DTOs and domain entities
- No framework dependencies in service layer
- Business logic calculations in service layer (not repository or handler)
- Use goroutines for concurrent processing in handlers
- Add toDomain() and toModel() converter functions in repository
```

**Usage:**
1. Copy template above
2. Replace `[CHANGE_THIS]` and `[MODULE_NAME]` with your module name
3. Replace table structure with your actual PostgreSQL table
4. Paste to AI assistant and generate
5. Copy generated files to your project
6. **Validate:** Run `go build ./internal/[MODULE_NAME]` to ensure no errors

**Post-Generation Checklist:**
- [ ] All files compile without errors
- [ ] Import paths use "github.com/golang-fiber-jwt"
- [ ] Handler methods use manual DTO parsing (not generic helpers)
- [ ] Repository has toDomain() and toModel() converter functions
- [ ] Service layer has no framework imports
- [ ] UUID fields handle nil values properly

---

**📚 For detailed code examples and complete guide, see:**

→ **[Step-by-Step Module Creation Guide](./docs/clean-architecture.md#quick-reference-new-module-checklist)**

---

### Testing

The project follows a two-tier testing strategy:

**Unit Tests** - Test business logic with mocked dependencies
```bash
# Run domain unit tests
go test ./internal/auth/...
```

**Integration Tests** - Test HTTP handlers with mocked services  
```bash
# Run transport integration tests
go test ./transport/http/auth/...
```

**All Tests**
```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

**Test Structure:**
- `/internal/<domain>/service_test.go` - Unit tests for business logic
- `/transport/http/<domain>/handler_integration_test.go` - Integration tests for HTTP layer

---

### Code Quality Checks

**Verify clean architecture compliance:**

```bash
# Build should succeed for each module
go build ./internal/user
go build ./internal/auth

# Domain entities should have no JSON tags
grep -r "json:" internal/*/entity.go  # Should be empty or only in database models

# Service layer should have no framework imports
grep -r "gorm\|fiber" internal/*/service.go  # Should be empty

# Check for correct import paths
grep -r "github.com/golang-fiber-jwt" internal/*/  # Should match module name

# No circular imports
go build -v ./...
```

---

### Architecture Checklist

When adding new code, verify:

- [ ] Domain entities have NO JSON/validation tags (only in DTOs)
- [ ] Service layer has NO framework imports (gorm, fiber, etc.)
- [ ] Business logic calculations are in service layer (not repository/handler)
- [ ] Repository interface and implementation in same file
- [ ] DTOs have proper JSON and validation tags
- [ ] Handler uses manual DTO mapping (no generic helpers)
- [ ] Handler implements concurrency where beneficial (goroutines + channels)
- [ ] Import paths use "github.com/golang-fiber-jwt" module name
- [ ] UUID fields handle nil values with proper type conversion
- [ ] Soft delete uses *time.Time for DeletedAt field
- [ ] All files compile: `go build ./internal/[module]`

---

### Detailed Documentation

For comprehensive guides on:
- **Creating new modules** with code examples
- **Testing strategy** (unit + integration tests)
- **Architecture patterns** and best practices
- **Domain examples** (Product, Order, etc.)

📖 **[Read the Complete Architecture Guide →](./docs/clean-architecture.md)**

## Contributing

1. Follow clean architecture principles
2. Keep domain layer pure (no framework imports)
3. Use shared utilities in `/pkg/`
4. Add tests for business logic
5. Document public APIs

## License

MIT



