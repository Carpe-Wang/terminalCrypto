.PHONY: build install clean test run help

# Build variables
BINARY_NAME=terminalcrypto
BUILD_DIR=bin
MAIN_FILE=main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the project
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated!"

# Run the application (for development)
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	@echo "Multi-platform build complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  make build      - Build the binary"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Download and tidy dependencies"
	@echo "  make run        - Build and run the application"
	@echo "  make build-all  - Build for multiple platforms"
	@echo "  make help       - Show this help message"
