.PHONY: help build run test clean docker-build docker-up docker-down docker-restart docker-logs docker-clean swagger deps fmt lint install-tools

# Variables
APP_NAME=minos
MAIN_PATH=./cmd/main.go
BINARY_NAME=tmp/$(APP_NAME)
DOCKER_COMPOSE=docker-compose
SWAGGER_CMD=swag

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  make build          - Build the application binary"
	@echo "  make run            - Run the application locally"
	@echo "  make dev            - Run with live reload (Air)"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-up      - Start all services with Docker Compose"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-restart - Restart all services"
	@echo "  make docker-logs    - View logs from all services"
	@echo "  make docker-clean   - Stop services and remove volumes"
	@echo ""
	@echo "Database commands:"
	@echo "  make db-shell       - Open PostgreSQL shell"
	@echo "  make db-migrate     - Run database migrations"
	@echo ""
	@echo "Development commands:"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make install-tools  - Install development tools"

## build: Build the application binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Build complete: $(BINARY_NAME)"

## run: Run the application locally (requires PostgreSQL to be running)
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

## dev: Run the application with live reload using Air
dev:
	@echo "Starting $(APP_NAME) with live reload..."
	@air -c .air.toml

## test: Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "✓ Tests complete"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

## clean: Clean build artifacts and cache
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/ tmp/
	@rm -f coverage.out coverage.html build-errors.log
	@go clean
	@echo "✓ Clean complete"

## docker-build: Build Docker images
docker-build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE) build
	@echo "✓ Docker build complete"

## docker-up: Start all services with Docker Compose
docker-up:
	@echo "Starting services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "✓ Services started"
	@echo "API available at: http://localhost:8080"
	@echo "Swagger docs at: http://localhost:8080/swagger/index.html"

## docker-down: Stop all services
docker-down:
	@echo "Stopping services..."
	@$(DOCKER_COMPOSE) down
	@echo "✓ Services stopped"

## docker-restart: Restart all services
docker-restart: docker-down docker-up

## docker-logs: View logs from all services
docker-logs:
	@$(DOCKER_COMPOSE) logs -f

## docker-clean: Stop services and remove volumes
docker-clean:
	@echo "Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down -v
	@echo "✓ Docker cleanup complete"

## db-shell: Open PostgreSQL shell
db-shell:
	@echo "Opening PostgreSQL shell..."
	@$(DOCKER_COMPOSE) exec postgres psql -U postgres -d minos

## db-migrate: Run database migrations (auto-migrate on app start)
db-migrate:
	@echo "Migrations run automatically when the app starts"
	@echo "To manually trigger: make docker-restart"

## swagger: Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@$(SWAGGER_CMD) init -g cmd/main.go -o ./docs
	@echo "✓ Swagger docs generated"

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "✓ Linting complete"; \
	else \
		echo "golangci-lint not installed. Run 'make install-tools' to install it."; \
	fi

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Tools installed"

# Default target
.DEFAULT_GOAL := help

