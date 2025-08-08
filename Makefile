# MH-Go Makefile

.PHONY: build run dev test clean docker-up docker-down db-backup db-restore db-reset

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Development with hot reload
dev:
	air

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-reset:
	docker-compose down -v
	docker-compose up -d

# Database commands
db-backup:
	./db-utils.sh backup

db-restore:
	@echo "Usage: make db-restore FILE=<backup_file>"
	@if [ -n "$(FILE)" ]; then ./db-utils.sh restore $(FILE); fi

db-reset:
	./db-utils.sh reset

db-list:
	./db-utils.sh list-backups

# Setup development environment
setup:
	./setup-dev.sh

# Help
help:
	@echo "MH-Go Development Commands:"
	@echo ""
	@echo "Building & Running:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Build and run the application"
	@echo "  make dev         - Run with hot reload (air)"
	@echo "  make test        - Run tests"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up   - Start all services"
	@echo "  make docker-down - Stop all services"
	@echo "  make docker-reset- Reset all services and volumes"
	@echo ""
	@echo "Database:"
	@echo "  make db-backup   - Create database backup"
	@echo "  make db-restore FILE=<file> - Restore from backup"
	@echo "  make db-reset    - Reset database (delete all data)"
	@echo "  make db-list     - List available backups"
	@echo ""
	@echo "Setup:"
	@echo "  make setup       - Setup development environment"