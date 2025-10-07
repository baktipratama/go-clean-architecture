.PHONY: build test run clean help

# Default environment variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= go_clean_code
DB_SSLMODE ?= disable
SERVER_PORT ?= 8081

# Build the application
build:
	@echo "Building the application..."
	go build -o bin/api cmd/api/*.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run the server
run:
	@echo "Starting the server on port $(SERVER_PORT)..."
	env SERVER_PORT=$(SERVER_PORT) DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_SSLMODE=$(DB_SSLMODE) go run cmd/api/*.go

# Run the built binary
run-binary: build
	@echo "Starting the server from binary on port $(SERVER_PORT)..."
	env SERVER_PORT=$(SERVER_PORT) DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_SSLMODE=$(DB_SSLMODE) ./bin/api

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Create binary directory
bin:
	mkdir -p bin

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env && echo "Created .env file from .env.example"; fi

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  run           - Run the server"
	@echo "  run-binary    - Build and run the binary"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  dev-setup     - Set up development environment"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  SERVER_PORT   - Server port (default: 8081)"
	@echo "  DB_HOST       - Database host (default: localhost)"
	@echo "  DB_PORT       - Database port (default: 5432)"
	@echo "  DB_USER       - Database user (default: postgres)"
	@echo "  DB_PASSWORD   - Database password (default: postgres)"
	@echo "  DB_NAME       - Database name (default: go_clean_code)"
	@echo "  DB_SSLMODE    - Database SSL mode (default: disable)"