#!/bin/bash

# MH-Go Development Script with Hot Reloading

echo "🚀 Starting MH-Go development environment..."

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Start PostgreSQL if not running
echo "📦 Starting PostgreSQL..."
docker-compose up postgres -d

# Wait for PostgreSQL
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Check if Air is installed
if ! command_exists air; then
    echo "📦 Installing Air for hot reloading..."
    go install github.com/air-verse/air@latest
fi

# Check if Templ is installed
if ! command_exists templ; then
    echo "📦 Installing Templ..."
    go install github.com/a-h/templ/cmd/templ@latest
fi

# Generate templ files
echo "🎨 Generating Templ templates..."
templ generate

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download

echo "✅ Development environment ready!"
echo ""
echo "🌟 Starting development server with hot reloading..."
echo "📝 Edit .go, .templ files and see changes instantly!"
echo "🌐 Application will be available at: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Start the development server with Air
if command -v air >/dev/null 2>&1; then
    air
elif [ -f "$HOME/go/bin/air" ]; then
    "$HOME/go/bin/air"
else
    echo "❌ Air not found. Installing Air..."
    go install github.com/air-verse/air@v1.61.1
    if [ -f "$HOME/go/bin/air" ]; then
        "$HOME/go/bin/air"
    else
        echo "❌ Failed to install Air. Please install manually: go install github.com/air-verse/air@v1.61.1"
        exit 1
    fi
fi
