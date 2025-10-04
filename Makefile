# Photo Video Organizer - Build Configuration
# Variables
BINARY_NAME := photo-organizer
BUILD_DIR := build
MAIN_PATH := cmd/organizer/main.go
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

# Detect OS for cross-platform commands
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RM := if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
    MKDIR := if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
    BINARY_EXT := .exe
else
    DETECTED_OS := $(shell uname -s)
    RM := rm -rf $(BUILD_DIR)
    MKDIR := mkdir -p $(BUILD_DIR)
    BINARY_EXT :=
endif

# Default target - build for all platforms
.PHONY: all
all: build-all

# Build the application
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME) v$(VERSION) for $(DETECTED_OS)..."
	@$(MKDIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT)"

# Run the application
.PHONY: run
run: build
	@echo "🚀 Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)$(BINARY_EXT)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@$(RM)
	@go clean
	@echo "✅ Clean complete"

# Run tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests complete"

# Run tests with coverage report
.PHONY: test-coverage
test-coverage: test
	@echo "📊 Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Install dependencies
.PHONY: deps
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@go mod verify
	@echo "✅ Dependencies installed"

# Format code
.PHONY: fmt
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"



# Cross-platform builds
.PHONY: build-all
build-all: clean
	@echo "🌍 Building for multiple platforms..."
	@$(MKDIR)
	@echo "  📱 Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "  🍎 Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "  🍎 Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "  🐧 Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "  🐧 Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "✅ Cross-platform builds complete"

# Development workflow
.PHONY: dev
dev: clean deps fmt test build
	@echo "🎉 Development build ready!"

# Release workflow
.PHONY: release
release: clean deps fmt test build-all test-coverage
	@echo "🚀 Release build complete!"

# Install binary to system
.PHONY: install
install: build
	@echo "📥 Installing $(BINARY_NAME)..."
ifeq ($(DETECTED_OS),Windows)
	@copy $(BUILD_DIR)\$(BINARY_NAME).exe %GOPATH%\bin\
else
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
endif
	@echo "✅ $(BINARY_NAME) installed to $(GOPATH)/bin/"

# Uninstall binary from system
.PHONY: uninstall
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
ifeq ($(DETECTED_OS),Windows)
	@del %GOPATH%\bin\$(BINARY_NAME).exe 2>nul || echo "Binary not found"
else
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)
endif
	@echo "✅ $(BINARY_NAME) uninstalled"

# Show help
.PHONY: help
help:
	@echo "📋 Available targets:"
	@echo "  all        - Build for all platforms (default)"
	@echo "  build      - Build the application for current platform"
	@echo "  build-all  - Cross-platform builds"
	@echo "  run        - Build and run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  fmt        - Format code"
	@echo "  dev        - Development workflow (deps, fmt, test, build)"
	@echo "  release    - Release workflow (all + coverage)"
	@echo "  install    - Install binary to system"
	@echo "  uninstall  - Remove binary from system"
	@echo "  help       - Show this help"
