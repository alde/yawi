# YAWI Makefile - keeping builds simple and tidy

.PHONY: help build clean test test-verbose lint deps install check

# Default target
help:
	@echo "YAWI - Yet Another Window Inspector"
	@echo ""
	@echo "Available targets:"
	@echo "  build         Build the binary for current platform"
	@echo "  clean         Clean up build artifacts"
	@echo "  test          Run tests"
	@echo "  test-verbose  Run tests with verbose output"
	@echo "  lint          Run linters"
	@echo "  deps          Download dependencies" 
	@echo "  install       Install to GOPATH/bin"
	@echo "  check         Quick build and test"
	@echo ""

# Build for current platform
build:
	@echo "Building YAWI..."
	go build -ldflags "-s -w" -o yawi ./cmd/yawi
	@echo "Build complete! Binary: ./yawi"

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f yawi
	rm -rf dist/
	go clean
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "Tests complete!"

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...
	@echo "Tests complete!"

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi
	@echo "Linting complete!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated!"

# Install to GOPATH/bin
install:
	@echo "Installing YAWI..."
	go install ./cmd/yawi
	@echo "YAWI installed!"


# Quick build and test
check: build test
	@echo "Quick check complete!"