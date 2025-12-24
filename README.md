# Minos

Minos is an AI reasoning service designed to execute complex AI logic and orchestrate LLM calls. It provides structured output transformation, prompt template management, and ephemeral session context handling.

## Architecture

- **Stateless Design**: Pure function approach for AI logic execution
- **Data Sources**: PostgreSQL for prompt templates, Redis for short-lived session state
- **API**: RESTful endpoints with Swagger documentation
- **Framework**: Go with Gin web framework

## Development Guide

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make
- Air (Optional for live reload)

### Quick Start

```bash
# Clone repository and setup
go mod download

# Start development environment with Docker
make docker-up

# Run with live reload
make dev
```

### Available Commands

**Development:**
- `make dev` - Run with live reload using Air
- `make run` - Run locally (requires PostgreSQL)
- `make build` - Build binary
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage

**Docker:**
- `make docker-up` - Start all services
- `make docker-down` - Stop services
- `make docker-logs` - View logs
- `make docker-restart` - Restart services

**Code Quality:**
- `make fmt` - Format code
- `make lint` - Run linter
- `make swagger` - Generate API documentation

**Database:**
- `make db-shell` - Open PostgreSQL shell
- `make db-migrate` - Run migrations

Run `make help` for all available commands.

### API Documentation

When running, access Swagger documentation at: http://localhost:8080/swagger/index.html