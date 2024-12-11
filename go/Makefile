# Project variables
BINARY_NAME=supavault
SOURCE_FILE=./cmd/supavault/main.go
BUILD_DIR=./bin
VERSION := $(or $(VERSION), 0.0.0)

# Detect GOOS and GOARCH from the system
GOOS := $(or $(GOOS), $(shell go env GOOS))
GOARCH := $(or $(GOARCH), $(shell go env GOARCH))

# Default target: Build the binary
all: build

# Build the application
build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)_$(GOOS)_$(GOARCH)_$(VERSION) \
		-ldflags="-X github.com/inspire-labs-tms-tech/supavault/pkg.Version=$(VERSION)" \
		$(SOURCE_FILE)

# Run the application
run: build
	@echo "Running the application..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Format the code
format:
	@echo "Formatting the code..."
	go fmt ./...

# Lint the code (requires golangci-lint)
lint:
	@echo "Linting the code..."
	golangci-lint run

# Tidy up Go modules
tidy:
	@echo "Tidying up Go modules..."
	go mod tidy

# Clean up generated files
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

# Cross-compile for a specific platform
build-linux:
	@echo "Cross-compiling for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(SOURCE_FILE)

build-windows:
	@echo "Cross-compiling for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(SOURCE_FILE)

# Help
help:
	@echo "Makefile for Go project:"
	@echo "  make all          - Build the application"
	@echo "  make run          - Build and run the application"
	@echo "  make test         - Run tests"
	@echo "  make format       - Format Go code"
	@echo "  make lint         - Run linters (requires golangci-lint)"
	@echo "  make tidy         - Tidy up Go modules"
	@echo "  make clean        - Remove build files"
	@echo "  make build-linux  - Cross-compile for Linux"
	@echo "  make build-windows- Cross-compile for Windows"

# Phony targets (not associated with actual files)
.PHONY: all build run test format lint tidy clean build-linux build-windows help