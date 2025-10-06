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

## 🚀 Features

This example demonstrates:

- **Clean Architecture**: Separation of concerns with clear dependency boundaries
- **User Management**: Complete CRUD operations for user entities  
- **PostgreSQL Support**: Production-ready database with migrations
- **RESTful API**: HTTP endpoints following REST conventions
- **Dependency Injection**: Clean dependency management with container pattern
- **Database Migrations**: Automated schema management
- **Environment Configuration**: Flexible configuration via environment variables
- **Testing Examples**: Unit tests showing how to test each layer
- **Repository Pattern**: Interface-based data access layer

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

## 🐳 Docker Support

*Coming soon - Docker containerization for easy deployment*

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

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