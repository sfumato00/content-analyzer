.PHONY: help install test build run docker-up docker-down docker-logs docker-rebuild clean lint fmt migrate-up migrate-down migrate-create verify

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
install: ## Install Go dependencies
	cd backend && go mod download && go mod tidy

test: ## Run all tests
	cd backend && go test ./... -v

test-coverage: ## Run tests with coverage
	cd backend && go test ./... -coverprofile=coverage.out
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

lint: ## Run linter
	cd backend && go vet ./...
	cd backend && gofmt -l .

fmt: ## Format Go code
	cd backend && go fmt ./...

# Build
build: ## Build the backend binary
	cd backend && go build -o ../bin/api cmd/api/main.go
	@echo "Binary built: bin/api"

build-linux: ## Build for Linux (useful for Docker)
	cd backend && GOOS=linux GOARCH=amd64 go build -o ../bin/api-linux cmd/api/main.go

# Run
run: ## Run the backend server (requires Docker services)
	cd backend && go run cmd/api/main.go

run-dev: docker-up run ## Start Docker services and run the server

# Docker
docker-up: ## Start all Docker services (postgres, redis, api)
	docker-compose up -d

docker-up-db: ## Start only database services (postgres, redis)
	docker-compose up -d postgres redis

docker-down: ## Stop all Docker services
	docker-compose down

docker-logs: ## Show logs from all services
	docker-compose logs -f

docker-logs-api: ## Show logs from API service only
	docker-compose logs -f api

docker-rebuild: ## Rebuild and restart all services
	docker-compose up -d --build

docker-clean: ## Remove all containers, volumes, and images
	docker-compose down -v
	docker system prune -f

# Database migrations
migrate-up: ## Run all pending migrations
	cd backend && go run cmd/api/main.go migrate up

migrate-down: ## Rollback last migration
	cd backend && go run cmd/api/main.go migrate down

migrate-create: ## Create a new migration (usage: make migrate-create name=add_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter required. Usage: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	cd backend && migrate create -ext sql -dir migrations -seq $(name)

# Verification
verify: ## Verify the setup is working
	@echo "Verifying setup..."
	@echo "\n1. Checking Docker services..."
	@docker-compose ps
	@echo "\n2. Checking database connection..."
	@docker-compose exec -T postgres pg_isready -U postgres || echo "❌ PostgreSQL not ready"
	@echo "\n3. Checking Redis connection..."
	@docker-compose exec -T redis redis-cli ping || echo "❌ Redis not ready"
	@echo "\n4. Checking API health..."
	@curl -s http://localhost:8080/health | jq . || echo "❌ API not responding"
	@echo "\n✅ Verification complete"

health: ## Check API health endpoint
	@curl -s http://localhost:8080/health | jq .

# Database access
db-shell: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d content_analyzer

db-dump: ## Dump database to SQL file
	docker-compose exec -T postgres pg_dump -U postgres content_analyzer > backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Database dumped to backup_*.sql"

db-restore: ## Restore database from SQL file (usage: make db-restore file=backup.sql)
	@if [ -z "$(file)" ]; then \
		echo "Error: file parameter required. Usage: make db-restore file=backup.sql"; \
		exit 1; \
	fi
	docker-compose exec -T postgres psql -U postgres -d content_analyzer < $(file)

redis-cli: ## Open Redis CLI
	docker-compose exec redis redis-cli

# Clean up
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf backend/coverage.out backend/coverage.html
	rm -f backend/*.log

clean-all: clean docker-clean ## Clean everything including Docker volumes

# Testing API endpoints
test-register: ## Test user registration endpoint
	@curl -X POST http://localhost:8080/api/v1/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' | jq .

test-login: ## Test user login endpoint
	@curl -X POST http://localhost:8080/api/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"password123"}' | jq .

# Development helpers
dev-setup: install docker-up-db ## Complete development setup
	@echo "✅ Development environment ready!"
	@echo "Run 'make run' to start the API server"

# CI/CD
ci: lint test ## Run CI checks (lint + test)

# Production
prod-build: ## Build for production
	cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/api cmd/api/main.go
