# ==================================
# issue2md Project Makefile
# GitHub Issue to Markdown Converter
# ==================================

# Project configuration
BINARY_CLI_NAME = issue2md-cli
BINARY_WEB_NAME = issue2md-web
DOCKER_IMAGE_NAME = issue2md
DOCKER_TAG ?= latest
BUILD_DIR = bin
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_SHA = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go configuration
GO_VERSION = 1.21
GO_OS = $(shell go env GOOS)
GO_ARCH = $(shell go env GOARCH)
LDFLAGS = -ldflags "-w -s -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commitSHA=$(COMMIT_SHA)"

# Docker configuration
DOCKER_BUILDKIT = 1
DOCKER_BUILD_ARGS = --build-arg VERSION=$(VERSION) --build-arg BUILD_TIME="$(BUILD_TIME)" --build-arg COMMIT_SHA=$(COMMIT_SHA)

# Default target
.DEFAULT_GOAL := help

# ==================================
# Core Targets
# ==================================

.PHONY: all build test lint docker-build clean

# Build all applications
all: build

# Build both CLI and Web applications
build:
	@echo "Building $(BINARY_CLI_NAME) and $(BINARY_WEB_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building CLI application..."
	CGO_ENABLED=0 GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_CLI_NAME) ./cmd/issue2md
	@echo "Building Web application..."
	CGO_ENABLED=0 GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_WEB_NAME) ./cmd/issue2mdweb
	@echo "✓ Build completed successfully!"
	@echo "  CLI:  $(BUILD_DIR)/$(BINARY_CLI_NAME)"
	@echo "  Web:  $(BUILD_DIR)/$(BINARY_WEB_NAME)"

# Run all unit tests
test:
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ All tests passed!"

# Run linter for static analysis
lint:
	@echo "Running static analysis with golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "⚠️  golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	golangci-lint run --timeout=5m
	@echo "✓ Linting completed successfully!"

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)..."
	@echo "Build args: VERSION=$(VERSION), BUILD_TIME=$(BUILD_TIME), COMMIT_SHA=$(COMMIT_SHA)"
	DOCKER_BUILDKIT=$(DOCKER_BUILDKIT) docker build $(DOCKER_BUILD_ARGS) -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .
	@echo "✓ Docker image built successfully!"
	@docker images | grep $(DOCKER_IMAGE_NAME)

# Clean all build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -testcache -modcache
	@docker system prune -f 2>/dev/null || true
	@echo "✓ Cleanup completed!"

# ==================================
# Additional Utility Targets
# ==================================

.PHONY: help test-coverage test-benchmark install docker-run docker-clean dev-setup format update-deps verify

# Show help
help:
	@echo "issue2md - GitHub Issue to Markdown Converter"
	@echo ""
	@echo "Core Targets:"
	@echo "  build          Build CLI and Web applications"
	@echo "  test           Run unit tests"
	@echo "  lint           Run static code analysis"
	@echo "  docker-build   Build Docker image"
	@echo "  clean          Clean build artifacts"
	@echo ""
	@echo "Additional Targets:"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  test-benchmark Run benchmark tests"
	@echo "  format         Format Go code"
	@echo "  install        Install applications to GOPATH/bin"
	@echo "  docker-run     Run Docker container"
	@echo "  docker-clean   Remove Docker image"
	@echo "  dev-setup      Setup development environment"
	@echo "  update-deps    Update Go dependencies"
	@echo "  verify         Run full CI pipeline (format+lint+test)"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION        $(VERSION)"
	@echo "  BUILD_TIME     $(BUILD_TIME)"
	@echo "  COMMIT_SHA     $(COMMIT_SHA)"
	@echo "  GO_OS          $(GO_OS)"
	@echo "  GO_ARCH        $(GO_ARCH)"

# Run tests with coverage
test-coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep -E "^total:" | awk '{print "Total coverage: " $$3}'

# Run benchmark tests
test-benchmark:
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem ./...

# Format Go code
format:
	@echo "Formatting Go code..."
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "⚠️  goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi
	@echo "✓ Code formatted successfully!"

# Install applications
install: build
	@echo "Installing applications to GOPATH/bin..."
	@mkdir -p $$(go env GOPATH)/bin
	@cp $(BUILD_DIR)/$(BINARY_CLI_NAME) $$(go env GOPATH)/bin/
	@cp $(BUILD_DIR)/$(BINARY_WEB_NAME) $$(go env GOPATH)/bin/
	@echo "✓ Applications installed successfully!"

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run --rm -it -p 8080:8080 -e PORT=8080 -e GITHUB_TOKEN=$(GITHUB_TOKEN) $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Clean Docker resources
docker-clean:
	@echo "Cleaning Docker resources..."
	@docker rmi -f $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) 2>/dev/null || true
	@echo "✓ Docker cleanup completed!"

# Setup development environment
dev-setup:
	@echo "Setting up development environment..."
	@go version
	@if [ $$(go version | awk '{print $$3}' | cut -d. -f2) -lt $(GO_VERSION) ]; then \
		echo "⚠️  Go version < $(GO_VERSION). Please upgrade Go."; \
	fi
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "✓ Development environment setup completed!"

# Update Go dependencies
update-deps:
	@echo "Updating Go dependencies..."
	@go get -u ./...
	@go mod tidy
	@go mod verify
	@echo "✓ Dependencies updated successfully!"

# Run verification (format + lint + test)
verify: format lint test
	@echo "✓ All verification checks passed!"