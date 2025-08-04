#!/bin/bash

# MH-Go Development Setup Script

echo "🚀 Setting up MH-Go development environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start PostgreSQL
echo "📦 Starting PostgreSQL container..."
docker-compose up postgres -d

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ and try again."
    exit 1
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download

# Check if Air is installed
if ! command -v air &> /dev/null; then
    echo "📦 Installing Air for hot reloading..."
    go install github.com/air-verse/air@latest
fi

# Run migrations
echo "🗄️ Running database migrations..."
go run cmd/server/main.go &
SERVER_PID=$!
sleep 2
kill $SERVER_PID 2>/dev/null || true

echo "✅ Development environment is ready!"
echo ""
echo "🌟 Quick Start Commands:"
echo "  • Start development server: air"
echo "  • Build application: go build cmd/server/main.go"
echo "  • Run tests: go test ./..."
echo "  • Start with Docker: docker-compose up"
echo ""
echo "🌐 Application will be available at: http://localhost:8080"
echo ""
echo "Happy coding! 🎉"
