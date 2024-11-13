# Makefile

# Go environment settings
GO=go
GINKGO=$(shell go env GOPATH)/bin/ginkgo
BINARY_NAME=cli

# Directories
SRC_DIR=.
BUILD_DIR=bin
E2E_TEST_DIR=e2e

# Variables
VERSION=1.0.0

# Build the CLI tool
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	$(GO) test -v ./...
	@echo "Unit tests completed."

# Run E2E tests with Ginkgo
.PHONY: e2e
e2e: build ginkgo
	@echo "Running E2E tests..."
	$(GINKGO) -v $(E2E_TEST_DIR)
	@echo "E2E tests completed."

# Install Ginkgo if not already installed
.PHONY: ginkgo
ginkgo:
ifeq (, $(shell which ginkgo))
	@echo "Installing Ginkgo..."
	$(GO) install github.com/onsi/ginkgo/v2/ginkgo@latest
endif

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	$(GO) clean
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Cleanup complete."

# Run all (build, test, e2e)
.PHONY: all
all: clean build test e2e

# Show help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build      - Build the CLI tool"
	@echo "  test       - Run unit tests"
	@echo "  e2e        - Run end-to-end (E2E) tests"
	@echo "  clean      - Clean up build artifacts"
	@echo "  all        - Run build, test, and e2e"
	@echo "  help       - Show this help message"