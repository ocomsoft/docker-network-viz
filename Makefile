# Makefile for docker-network-viz
# Docker CLI plugin for visualizing Docker network topology

# Binary name
BINARY_NAME := docker-network-viz

# Build directories
BUILD_DIR := bin
DIST_DIR := dist

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOFMT := goimports

# Main package path
MAIN_PKG := ./cmd/docker-network-viz

# Installation directories
USER_PLUGIN_DIR := $(HOME)/.docker/cli-plugins
SYSTEM_PLUGIN_DIR := /usr/local/lib/docker/cli-plugins

# Build flags
LDFLAGS := -s -w
BUILD_FLAGS := -ldflags "$(LDFLAGS)"

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: all build test lint clean install install-system help fmt tidy deps

## help: Show this help message
help:
	@echo "docker-network-viz - Docker CLI plugin for network visualization"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^##[[:space:]]/ { \
		split($$0, parts, ": "); \
		target = substr(parts[1], 4); \
		desc = parts[2]; \
		for (i = 3; i <= length(parts); i++) desc = desc ": " parts[i]; \
		printf "  %-20s %s\n", target, desc; \
	}' $(MAKEFILE_LIST)

## all: Build and test the project
all: lint test build

## build: Compile the binary to docker-network-viz
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PKG)
	@echo "Build complete: $(BINARY_NAME)"

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -cover ./...

## lint: Run golangci-lint and goimports
lint: fmt
	@echo "Running golangci-lint..."
	golangci-lint run -v ./...
	@echo "Lint complete."

## fmt: Format code with goimports
fmt:
	@echo "Formatting code with goimports..."
	$(GOFMT) -w .
	@echo "Formatting complete."

## install: Install the binary to ~/.docker/cli-plugins/
install: build
	@echo "Installing $(BINARY_NAME) to $(USER_PLUGIN_DIR)..."
	@mkdir -p $(USER_PLUGIN_DIR)
	@cp $(BINARY_NAME) $(USER_PLUGIN_DIR)/$(BINARY_NAME)
	@chmod +x $(USER_PLUGIN_DIR)/$(BINARY_NAME)
	@echo "Installation complete."
	@echo "You can now use: docker network-viz"

## install-system: Install to /usr/local/lib/docker/cli-plugins/ (requires sudo)
install-system: build
	@echo "Installing $(BINARY_NAME) to $(SYSTEM_PLUGIN_DIR) (requires sudo)..."
	sudo mkdir -p $(SYSTEM_PLUGIN_DIR)
	sudo cp $(BINARY_NAME) $(SYSTEM_PLUGIN_DIR)/$(BINARY_NAME)
	sudo chmod +x $(SYSTEM_PLUGIN_DIR)/$(BINARY_NAME)
	@echo "System installation complete."
	@echo "You can now use: docker network-viz"

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete."

## tidy: Tidy and verify go modules
tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Modules tidy complete."

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "Dependencies downloaded."

## cover: Run tests with coverage report
cover:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## check: Run all checks (fmt, lint, test)
check: fmt lint test
	@echo "All checks passed."
