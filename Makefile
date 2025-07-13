# Helm Chart Browser - Makefile

# Build automation for cross-platform releases

# Application name

APP_NAME := helm-browser

# Version (can be overridden)

VERSION ?= $(shell git describe –tags –always –dirty 2>/dev/null || echo “dev”)

# Build flags

LDFLAGS := -s -w -X main.version=$(VERSION)
BUILD_FLAGS := -ldflags=”$(LDFLAGS)”

# Output directory

DIST_DIR := dist

# Default target

.PHONY: all
all: build

# Clean build artifacts

.PHONY: clean
clean:
@echo “🧹 Cleaning build artifacts…”
rm -rf $(DIST_DIR)
rm -f $(APP_NAME)
go clean

# Download dependencies

.PHONY: deps
deps:
@echo “📦 Downloading dependencies…”
go mod download
go mod tidy

# Run tests

.PHONY: test
test:
@echo “🧪 Running tests…”
go test -v ./…

# Run tests with coverage

.PHONY: test-coverage
test-coverage:
@echo “🧪 Running tests with coverage…”
go test -v -coverprofile=coverage.out ./…
go tool cover -html=coverage.out -o coverage.html
@echo “📊 Coverage report generated: coverage.html”

# Lint code

.PHONY: lint
lint:
@echo “🔍 Linting code…”
golangci-lint run

# Format code

.PHONY: fmt
fmt:
@echo “🎨 Formatting code…”
go fmt ./…

# Build for current platform

.PHONY: build
build: deps
@echo “🔨 Building $(APP_NAME) for current platform…”
go build $(BUILD_FLAGS) -o $(APP_NAME) .

# Build for all platforms

.PHONY: build-all
build-all: clean deps build-linux build-darwin build-windows
@echo “✅ All builds completed successfully!”

# Create dist directory

$(DIST_DIR):
mkdir -p $(DIST_DIR)

# Build for Linux AMD64

.PHONY: build-linux
build-linux: $(DIST_DIR)
@echo “🐧 Building for Linux AMD64…”
GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux .

# Build for Linux ARM64

.PHONY: build-linux-arm64
build-linux-arm64: $(DIST_DIR)
@echo “🐧 Building for Linux ARM64…”
GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .

# Build for macOS AMD64

.PHONY: build-darwin
build-darwin: $(DIST_DIR)
@echo “🍎 Building for macOS AMD64…”
GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .

# Build for macOS ARM64 (Apple Silicon)

.PHONY: build-darwin-arm64
build-darwin-arm64: $(DIST_DIR)
@echo “🍎 Building for macOS ARM64 (Apple Silicon)…”
GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .

# Build for Windows AMD64

.PHONY: build-windows
build-windows: $(DIST_DIR)
@echo “🪟 Building for Windows AMD64…”
GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows.exe .

# Build all platforms including ARM64

.PHONY: build-all-platforms
build-all-platforms: clean deps build-linux build-linux-arm64 build-darwin build-darwin-arm64 build-windows
@echo “✅ All platform builds completed successfully!”

# Run the application

.PHONY: run
run: build
@echo “🚀 Running $(APP_NAME)…”
./$(APP_NAME)

# Development mode with auto-rebuild

.PHONY: dev
dev:
@echo “🔄 Starting development mode…”
go run .

# Install dependencies for development

.PHONY: dev-deps
dev-deps:
@echo “🛠️  Installing development dependencies…”
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Package binaries with checksums

.PHONY: package
package: build-all-platforms
@echo “📦 Generating checksums…”
cd $(DIST_DIR) && sha256sum * > checksums.txt
@echo “📦 Packaging completed!”

# Release preparation

.PHONY: release-prep
release-prep: test lint package
@echo “🎉 Release preparation completed!”
@echo “📁 Binaries are in $(DIST_DIR)/”
@ls -la $(DIST_DIR)/

# Quick development check

.PHONY: check
check: fmt lint test
@echo “✅ All checks passed!”

# Show help

.PHONY: help
help:
@echo “🚀 Helm Chart Browser - Build Commands”
@echo “”
@echo “Development:”
@echo “  make dev          - Run in development mode”
@echo “  make run          - Build and run application”
@echo “  make test         - Run tests”
@echo “  make test-coverage- Run tests with coverage report”
@echo “  make fmt          - Format code”
@echo “  make lint         - Lint code”
@echo “  make check        - Run fmt, lint, and test”
@echo “”
@echo “Building:”
@echo “  make build        - Build for current platform”
@echo “  make build-all    - Build for Linux, macOS, Windows (AMD64)”
@echo “  make build-all-platforms - Build for all platforms including ARM64”
@echo “”
@echo “Platform-specific builds:”
@echo “  make build-linux  - Build for Linux AMD64”
@echo “  make build-linux-arm64 - Build for Linux ARM64”
@echo “  make build-darwin - Build for macOS AMD64”
@echo “  make build-darwin-arm64 - Build for macOS ARM64”
@echo “  make build-windows- Build for Windows AMD64”
@echo “”
@echo “Release:”
@echo “  make package      - Build all platforms and generate checksums”
@echo “  make release-prep - Full release preparation (test, lint, package)”
@echo “”
@echo “Utilities:”
@echo “  make deps         - Download dependencies”
@echo “  make dev-deps     - Install development dependencies”
@echo “  make clean        - Clean build artifacts”
@echo “  make help         - Show this help message”

# Default help target

.DEFAULT_GOAL := help