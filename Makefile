# Variables
APP_NAME = goVault

# Default target: help
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Build the Go project
build: ## Build the Go project
	@go get -u ./...
	@go build -o .build/$(APP_NAME)

# Run the Go project
run: build ## Build and run the Go project
	./.build/$(APP_NAME)

# Start Docker containers using Docker Compose
up: ## Start Docker containers
	docker compose build
	docker compose pull
	docker compose up -d

# Stop Docker containers using Docker Compose
down: ## Stop Docker containers
	@docker compose down


# Remove Docker containers and images
remove: ## Remove Docker containers and images
	@docker compose down --rmi all

# Clean the repository
clean: ## Clean build artifacts
	@go clean
	@rm -f $(APP_NAME)

# Run tests
test: ## Run tests
	@go test ./...

.PHONY: build run docker-start docker-stop docker-restart docker-remove clean test lint
