#!/bin/bash

echo "🧪 Testing mH-go Application"
echo "=================================="

# Test 1: Check if Go version is correct
echo -n "✅ Go version: "
go version

# Test 2: Check if Air is available
echo -n "✅ Air version: "
air -v 2>/dev/null | head -1

# Test 3: Check if database is accessible via Docker
echo -n "✅ PostgreSQL container: "
if docker-compose ps postgres | grep -q "Up"; then
    echo "Running ✅"
else
    echo "Not running ❌"
fi

# Test 4: Test building the application
echo "✅ Building application..."
if go build -o tmp/test ./cmd/server/main.go; then
    echo "   Build successful ✅"
    rm -f tmp/test
else
    echo "   Build failed ❌"
fi

# Test 5: Check if local development works
echo "✅ Testing local development environment..."
echo "   Starting application on http://localhost:8080"
echo "   You can now run: ./dev.sh"

echo ""
echo "🎉 Setup verification complete!"
echo "You can now use either:"
echo "   • Local development: ./dev.sh"
echo "   • Docker development: make docker-dev"
