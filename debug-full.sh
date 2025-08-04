#!/bin/bash

echo "🔍 Debugging Registration Issue"
echo "==============================="

# Check if PostgreSQL is running
echo "1. Checking PostgreSQL..."
if lsof -i :5432 > /dev/null 2>&1; then
    echo "   ✅ PostgreSQL is running on port 5432"
else
    echo "   ❌ PostgreSQL is NOT running on port 5432"
    echo "   Starting PostgreSQL with Docker..."

    # Clean up any existing container
    docker rm -f mh-go-postgres 2>/dev/null

    # Start PostgreSQL
    docker run -d --name mh-go-postgres \
        -e POSTGRES_DB=mh_go \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=password \
        -p 5432:5432 \
        postgres:15-alpine

    echo "   Waiting for PostgreSQL to start..."
    sleep 15

    if lsof -i :5432 > /dev/null 2>&1; then
        echo "   ✅ PostgreSQL started successfully"
    else
        echo "   ❌ Failed to start PostgreSQL"
        exit 1
    fi
fi

# Check if Go app compiles
echo "2. Checking Go compilation..."
if go build -o /tmp/mh-go-test ./cmd/server/main.go; then
    echo "   ✅ Go application compiles successfully"
    rm -f /tmp/mh-go-test
else
    echo "   ❌ Go compilation failed"
    exit 1
fi

# Start the application
echo "3. Starting the application..."
go run cmd/server/main.go &
APP_PID=$!

# Wait for app to start
sleep 5

# Check if app is running
echo "4. Checking if application is responding..."
if curl -s -o /dev/null -w "%{http_code}" http://localhost:8080 | grep -q "200"; then
    echo "   ✅ Application is responding"
else
    echo "   ❌ Application is not responding"
    kill $APP_PID 2>/dev/null
    exit 1
fi

# Test registration endpoint
echo "5. Testing registration endpoint..."
response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST http://localhost:8080/register \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "email=debug$(date +%s)@example.com&password=password123&confirm_password=password123")

echo "Response:"
echo "$response"

http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d: -f2)

if [ "$http_code" = "200" ]; then
    echo "   ✅ Registration successful"
else
    echo "   ❌ Registration failed with HTTP code: $http_code"
fi

# Clean up
echo "6. Cleaning up..."
kill $APP_PID 2>/dev/null
echo "   Application stopped"

echo ""
echo "🏁 Debug complete!"
