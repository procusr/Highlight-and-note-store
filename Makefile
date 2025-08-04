.PHONY: build run test clean install-deps dev docker-up docker-down migrate-up migrate-down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run with hot reloading using Air
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/

# Install dependencies
install-deps:
	go mod download
	go mod tidy

# Install development tools
install-tools:
	go install github.com/air-verse/air@v1.61.1
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/a-h/templ/cmd/templ@latest

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Development with Docker and hot reloading
docker-dev:
	docker-compose --profile dev up

# Production with Docker
docker-prod:
	docker-compose --profile prod up --build -d

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/mh_go?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/mh_go?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Setup development environment
setup-dev:
	make install-deps
	make install-tools
	make docker-up
	sleep 5
	make migrate-up

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate Tailwind CSS (if using build process)
css:
	npx tailwindcss -i ./static/input.css -o ./static/output.css --watch

# Generate Templ templates
templ:
	templ generate

# Watch and generate Templ templates
templ-watch:
	templ generate --watch
