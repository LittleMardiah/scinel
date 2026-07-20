.PHONY: help dev test test-coverage test-integration test-unit lint build build-docker db-up db-down db-migrate db-rollback db-seed db-reset docker-up docker-down docker-logs clean format

# Default target (show help)
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

help:
	@echo "ScIntel Development Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev              - Run backend with hot reload (requires air)"
	@echo "  make run              - Run backend directly (no hot reload)"
	@echo ""
	@echo "Testing:"
	@echo "  make test             - Run all tests (unit + integration)"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-integration - Run only integration tests"
	@echo "  make test-unit        - Run only unit tests"
	@echo "  make lint             - Run linter (golangci-lint)"
	@echo ""
	@echo "Build:"
	@echo "  make build            - Build binary (scinel)"
	@echo "  make build-docker     - Build Docker image"
	@echo ""
	@echo "Database:"
	@echo "  make db-up            - Start database containers"
	@echo "  make db-down          - Stop database containers"
	@echo "  make db-migrate       - Run all migrations"
	@echo "  make db-rollback      - Rollback last migration"
	@echo "  make db-seed          - Seed test data"
	@echo "  make db-reset         - Drop all tables & re-migrate"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up        - Start Docker Compose stack"
	@echo "  make docker-down      - Stop Docker Compose stack"
	@echo "  make docker-logs      - View Docker logs"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make format           - Format code (gofmt)"
	@echo ""

dev:
	@echo "Running with hot reload..."
	@which air > /dev/null || (echo "air not installed"; echo "Install: go install github.com/cosmtrek/air@latest"; exit 1)
	air

run:
	@echo "Running backend..."
	go run ./cmd/api

test:
	@echo "Running all tests..."
	go test -v -race ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-integration:
	@echo "Running integration tests..."
	go test -v -race -run Integration ./...

test-unit:
	@echo "Running unit tests..."
	go test -v -race -run Unit ./... -short

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed"; exit 1)
	golangci-lint run ./...

build:
	@echo "Building binary..."
	go build -o bin/scinel ./cmd/api
	@echo "Binary created: bin/scinel"

build-docker:
	@echo "Building Docker image..."
	docker build -t scinel:latest .
	@echo "Docker image created: scinel:latest"

db-up:
	@echo "Starting database services..."
	docker-compose up -d postgres redis
	@echo "Services started"
	@sleep 3
	@docker-compose ps

db-down:
	@echo "Stopping database services..."
	docker-compose down

db-migrate:
	@echo "Running migrations..."
	@which migrate > /dev/null || (echo "golang-migrate not installed"; echo "Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; exit 1)
	migrate -path migrations -database "}{DATABASE_URL}" up

db-rollback:
	@echo "Rolling back migration..."
	migrate -path migrations -database "}{DATABASE_URL}" down 1

db-seed:
	@echo "Seeding test data..."
	go run scripts/seed-db.go

db-reset:
	@echo "WARNING: This will drop all tables!"
	@read -p "Continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ }REPLY =~ ^[Yy]} ]]; then \
	migrate -path migrations -database "}{DATABASE_URL}" down -all; \
	migrate -path migrations -database "}{DATABASE_URL}" up; \
	echo "Database reset complete"; \
	fi

docker-up:
	@echo "Starting Docker Compose..."
	docker-compose up -d
	@echo "Services started"
	@sleep 2
	docker-compose ps

docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

docker-logs:
	@echo "Docker logs:"
	docker-compose logs -f

clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/scinel
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete"

format:
	@echo "Formatting code..."
	gofmt -s -w ./cmd ./internal ./pkg
	@echo "Format complete"