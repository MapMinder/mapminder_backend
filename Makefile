.PHONY: dev build run test lint clean install-air

# Development with hot reload
dev:
	@if ! command -v air > /dev/null 2>&1 && ! test -f $(shell go env GOPATH)/bin/air; then \
		echo "Air is not installed. Please install it first with: make install-air"; \
		exit 1; \
	fi
	@if command -v air > /dev/null 2>&1; then \
		air -c .air.toml; \
	else \
		$(shell go env GOPATH)/bin/air -c .air.toml; \
	fi

# Install air for hot reload
install-air:
	go install github.com/air-verse/air@latest

# Build the application
build:
	go build -o mapminder ./cmd/...

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test ./... -v

# Run specific test
test-health:
	go test ./feature/health/... -v

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available commands:"
	@echo "  dev           - Start development server with hot reload"
	@echo "  install-air   - Install air for hot reload"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-health   - Run health feature tests"
	@echo "  lint          - Run linter"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  help          - Show this help message"