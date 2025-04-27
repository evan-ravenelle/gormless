# Makefile for Go PostgreSQL Database Schema Manager

# Variables
BINARY_NAME=gormless
GO=go
GOCOVER=$(GO) tool cover
SRC=$(shell find . -name "*.go" -type f)
PKG=./...

# Default target
.PHONY: all
all: clean build test

# Build the application
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) -v

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Clean build artifacts
.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Run all tests
.PHONY: test
test:
	$(GO) test $(PKG) -v

# Run tests excluding integration tests
.PHONY: test-unit
test-unit:
	$(GO) test -short $(PKG) -v

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GOCOVER) -html=coverage.out

# Run integration tests only
.PHONY: test-integration
test-integration:
	$(GO) test -run "^_Test" -v $(PKG)

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt $(PKG)

# Verify dependencies
.PHONY: deps
deps:
	$(GO) mod verify
	$(GO) mod tidy

# Install required tools
.PHONY: tools
tools:
	$(GO) install github.com/stretchr/testify@latest
	$(GO) install github.com/DATA-DOG/go-sqlmock@latest

# Run tests in docker environment
.PHONY: docker-test
docker-test:
	docker run --rm -v $(PWD):/app -w /app golang:1.24 go test $(PKG)

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build           - Build the application"
	@echo "  make run             - Run the application"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make test            - Run all tests"
	@echo "  make test-unit       - Run unit tests only (no integration tests)"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make fmt             - Format code"
	@echo "  make deps            - Verify dependencies"
	@echo "  make tools           - Install required tools"
	@echo "  make docker-test     - Run tests in Docker environment"