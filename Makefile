.PHONY: help build test clean run docker-up docker-down install lint fmt

# Default target
help:
	@echo "Available targets:"
	@echo "  make build       - Build the wallet-tracker binary"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linters"
	@echo "  make fmt         - Format code"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make run         - Run the application"
	@echo "  make docker-up   - Start Docker services"
	@echo "  make docker-down - Stop Docker services"
	@echo "  make install     - Install the binary"

# Build variables
BINARY_NAME=wallet-tracker
BUILD_DIR=./build
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} cmd/wallet-tracker/main.go
	@echo "Build complete: ${BUILD_DIR}/${BINARY_NAME}"

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Test coverage:"
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}'

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		go vet ./...; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w ${GO_FILES}
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@rm -f coverage.out
	@echo "Clean complete"

# Run the application
run: build
	@echo "Running ${BINARY_NAME}..."
	${BUILD_DIR}/${BINARY_NAME}

# Docker commands
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@docker-compose ps

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "Services stopped"

# Install the binary
install: build
	@echo "Installing ${BINARY_NAME}..."
	go install cmd/wallet-tracker/main.go
	@echo "Installation complete"

# Development setup
dev-setup: docker-up
	@echo "Setting up development environment..."
	@cp .env.example .env 2>/dev/null || true
	@echo "Development environment ready!"
	@echo "Edit .env file with your configuration"

# Cross-compilation targets
build-all:
	@echo "Building for all platforms..."
	@mkdir -p ${BUILD_DIR}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 cmd/wallet-tracker/main.go
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 cmd/wallet-tracker/main.go
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 cmd/wallet-tracker/main.go
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 cmd/wallet-tracker/main.go
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe cmd/wallet-tracker/main.go
	@echo "Cross-compilation complete"
