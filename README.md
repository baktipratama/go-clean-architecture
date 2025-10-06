# Go Clean Code - User Management API

**Example implementation of a backend API using Clean Architecture principles in Go**

This is a demonstration project showcasing how to build a REST API with Clean Architecture patterns for user management with PostgreSQL database support. Perfect for learning and reference.

## 🏗️ Architecture

This example project demonstrates Clean Architecture implementation with the following layers:

```
cmd/api/           # Application entry point and configuration
├── main.go        # Application bootstrap
├── config.go      # Configuration management
├── database.go    # Database connection and migrations
├── container.go   # Dependency injection container
└── router.go      # HTTP route definitions

internal/
├── dto/           # Data Transfer Objects
├── handler/       # HTTP handlers (Controllers)
├── repository/    # Data access layer
└── usecase/       # Business logic layer

migrations/        # Database migrations
```

## 🏛️ Clean Architecture Implementation

This repository demonstrates how Clean Architecture principles are implemented in Go:

### 🔄 Dependency Rule
The core principle where dependencies point inward - outer layers depend on inner layers, never the reverse:

```
cmd/api (Framework & Drivers) 
    ↓ depends on
internal/handler (Interface Adapters)
    ↓ depends on  
internal/usecase (Application Business Rules)
    ↓ depends on
internal/repository (Interface - Enterprise Business Rules)
```

### 🎯 Layer Separation

**1. Data Layer (`internal/repository/`) - "Where we save and get data"**
- **Repository** = Data storage manager - handles saving/loading data from database
- Contains the `User` entity and `UserRepositoryInterface`
- Like a librarian who knows how to find and store books
- No dependencies on external frameworks or databases
- Pure business logic that could work in any context

**2. Business Logic Layer (`internal/usecase/`) - "What our app actually does"**
- **Usecase** = Business rules engine - implements what the application should do
- Contains `UserUsecase` implementing business workflows  
- Like a manager who makes decisions and coordinates work
- Orchestrates data flow between entities and repositories
- Depends only on repository interfaces, not implementations

**3. Presentation Layer (`internal/handler/`, `internal/dto/`) - "How users interact with our app"**
- **Handler** = Request processor - converts web requests into business operations
- **DTO** = Data messenger - carries data between different parts of the app
- Like a translator who speaks both "web language" and "business language"
- HTTP handlers transform web requests to use case inputs
- No business logic - pure data transformation

**4. Frameworks & Drivers (`cmd/api/`)**
- Contains framework-specific code (HTTP server, database connections)
- Configuration management and dependency injection
- Database migrations and connection setup
- Entry point that wires everything together

### 🔌 Dependency Inversion

**Interface Definitions**: Repository interfaces are defined in the business layer but implemented in the infrastructure layer:

```go
// Defined in internal/repository (inner layer)
type UserRepositoryInterface interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    // ...
}

// Implemented in internal/repository (same layer, different file)
type UserRepositoryImpl struct {
    db *sql.DB  // External dependency
}
```

**Dependency Injection**: Dependencies are injected from the outermost layer:

```go
// cmd/api/container.go
func NewContainer() *Container {
    // Infrastructure layer creates concrete implementations
    db := ConnectDatabase(&config.Database)
    userRepo := repository.NewUserRepository(db)
    
    // Business layer receives interfaces
    userUsecase := usecase.NewUserUsecase(userRepo)  // userRepo implements UserRepositoryInterface
    userHandler := handler.NewUserHandler(userUsecase)
    
    return &Container{
        UserHandler: userHandler,
    }
}
```

### 🧪 Testability

Each layer can be tested independently:

**Repository Tests**: Use SQL mocks to test data access logic without a real database
```go
func TestUserRepositoryImpl_Create(t *testing.T) {
    db, mock, err := sqlxmock.Newx()
    repo := &UserRepositoryImpl{db: db.DB}
    
    mock.ExpectExec("INSERT INTO users...").WillReturnResult(...)
    err := repo.Create(ctx, user)
    // Test database interaction logic
}
```

**Use Case Tests**: Mock repository interfaces to test business logic
```go
func TestUserUsecase_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    usecase := NewUserUsecase(mockRepo)
    
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    err := usecase.CreateUser(ctx, dto)
    // Test business logic without database
}
```

**Handler Tests**: Mock use cases to test HTTP handling
```go
func TestUserHandler_CreateUser(t *testing.T) {
    mockUsecase := &MockUserUsecase{}
    handler := NewUserHandler(mockUsecase)
    
    mockUsecase.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
    // Test HTTP request/response handling
}
```

### 📁 Directory Structure Benefits

```
internal/
├── dto/           # Data contracts (no business logic)
├── handler/       # HTTP adapters (framework concerns)
├── repository/    # Data access contracts & implementations
└── usecase/       # Pure business logic

cmd/api/           # Application composition & framework setup
├── main.go        # Entry point
├── config.go      # Configuration management
├── database.go    # Infrastructure setup
└── container.go   # Dependency injection
```

This structure ensures:
- **Screaming Architecture**: Directory names reveal what the application does, not what framework it uses
- **Plugin Architecture**: Database, web framework, etc. are plugins to the business logic
- **Framework Independence**: Business logic doesn't depend on Go-specific frameworks
- **Database Independence**: Can switch from PostgreSQL to MongoDB by changing only the repository implementation

## 📋 Prerequisites

- Go 1.24 or higher
- PostgreSQL (optional - falls back to in-memory storage)

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-clean-code
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Build the application**
   ```bash
   go build -o api cmd/api/*.go
   ```

## ⚙️ Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and adjust as needed:

```env
# Server Configuration
SERVER_PORT=8081

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_clean_code
DB_SSLMODE=disable
```

### Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8081` | HTTP server port |

| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `go_clean_code` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode for PostgreSQL |

## 🏃‍♂️ Running the Application

### Running with PostgreSQL
```bash
# Make sure PostgreSQL is running and accessible
go run cmd/api/*.go

# Or using the built binary
./api
```

The server will start on `http://localhost:8081` (or the port specified in `SERVER_PORT`).

## 📚 API Endpoints

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/users` | Get all users |
| `GET` | `/users/{id}` | Get user by ID |
| `POST` | `/users` | Create new user |
| `PUT` | `/users/{id}` | Update user |
| `DELETE` | `/users/{id}` | Delete user |

### Example Requests

**Create User:**
```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com"
  }'
```

**Get All Users:**
```bash
curl http://localhost:8081/users
```

**Get User by ID:**
```bash
curl http://localhost:8081/users/{user-id}
```

**Update User:**
```bash
curl -X PUT http://localhost:8081/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com"
  }'
```

**Delete User:**
```bash
curl -X DELETE http://localhost:8081/users/{user-id}
```

## 🗄️ Database

### PostgreSQL Setup

1. **Install PostgreSQL** (if not already installed)
2. **Create database:**
   ```sql
   CREATE DATABASE go_clean_code;
   ```
3. **Set environment variables** in `.env`:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=go_clean_code
   ```

### Migrations

Database migrations are automatically run on startup when using PostgreSQL. Migration files are located in the `migrations/` directory.

**Manual migration commands:**
```bash
# Up migrations
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up

# Down migrations
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" down
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test files
go test ./internal/handler/
go test ./internal/usecase/
go test ./internal/repository/
```

## 🔨 Development

### Project Structure Explained

- **`cmd/api/`**: Application-specific code (main, config, dependency injection)
- **`internal/dto/`**: Data Transfer Objects for API communication
- **`internal/handler/`**: HTTP handlers that handle requests/responses
- **`internal/usecase/`**: Business logic and use cases
- **`internal/repository/`**: Data access layer with interfaces and implementations
- **`migrations/`**: SQL migration files for database schema

### Adding New Features

1. **Add DTO** in `internal/dto/` if needed
2. **Create Repository Interface** in `internal/repository/`
3. **Implement Repository** for both memory and PostgreSQL
4. **Add Use Case** in `internal/usecase/`
5. **Create Handler** in `internal/handler/`
6. **Add Routes** in `cmd/api/router.go`
7. **Update Container** in `cmd/api/container.go` for dependency injection

## 🔍 Clean Architecture Benefits

This example implementation demonstrates:

- **Independence**: Business rules don't depend on external frameworks
- **Testability**: Business rules can be tested without UI, database, or external elements
- **Independence of UI**: UI can change without changing business rules
- **Independence of Database**: Business rules are not bound to the database
- **Independence of External Agency**: Business rules don't know about the outside world

## 🛠️ Built With

- [Go 1.24](https://golang.org/) - Programming language
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP router
- [PostgreSQL](https://www.postgresql.org/) - Database
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver
- [testify](https://github.com/stretchr/testify) - Testing toolkit