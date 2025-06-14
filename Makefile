.PHONY: help build run clean test deps migrate docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download and install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/building-management-system main.go

run: ## Run the application
	go run main.go

run-dev: ## Run with hot reload (requires air)
	air

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

migrate: ## Run database migrations manually
	@echo "Creating database if not exists..."
	-createdb building_management 2>/dev/null || true
	@echo "Running migrations..."
	go run main.go --migrate-only

db-setup: ## Setup local PostgreSQL database
	createdb building_management || true
	psql -d building_management -c "SELECT 1;" > /dev/null 2>&1 && echo "Database connection successful"

db-reset: ## Reset database (WARNING: This will drop all data)
	dropdb building_management --if-exists
	createdb building_management

seed: ## Seed database with sample data
	@echo "Seeding database..."
	curl -X POST http://localhost:3000/buildings \
		-H "Content-Type: application/json" \
		-d '{"name": "Sunrise Tower", "address": "123 Main Street, Downtown"}'
	curl -X POST http://localhost:3000/buildings \
		-H "Content-Type: application/json" \
		-d '{"name": "Moonlight Plaza", "address": "456 Oak Avenue, Uptown"}'
	curl -X POST http://localhost:3000/apartments \
		-H "Content-Type: application/json" \
		-d '{"building_id": 1, "number": "101", "floor": 1, "sq_meters": 75}'
	curl -X POST http://localhost:3000/apartments \
		-H "Content-Type: application/json" \
		-d '{"building_id": 1, "number": "102", "floor": 1, "sq_meters": 82}'

docker-build: ## Build Docker image
	docker build -t building-management-system .

docker-run: ## Run with Docker Compose
	docker-compose up -d

docker-stop: ## Stop Docker containers
	docker-compose down

sqlboiler-gen: ## Generate SQLBoiler models
	sqlboiler psql

install-dev-tools: ## Install development tools
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Default target
all: deps build

# Production targets
.PHONY: build-prod deploy

build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/building-management-system main.go 