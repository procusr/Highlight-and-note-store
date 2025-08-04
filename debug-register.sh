#!/bin/bash

echo "Starting application with debug output..."

# Start PostgreSQL if not running
if ! docker ps | grep -q postgres; then
    echo "Starting PostgreSQL container..."
    docker run -d --name mh-go-postgres \
        -e POSTGRES_DB=mh_go \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=password \
        -p 5432:5432 \
        postgres:15-alpine
    echo "Waiting for PostgreSQL to start..."
    sleep 10
fi

# Start the application
echo "Starting Go application..."
go run cmd/server/main.go &
APP_PID=$!

# Wait for application to start
sleep 3

# Test registration
echo "Testing registration..."
response=$(curl -v -X POST http://localhost:8080/register \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "email=test@example.com&password=password123&confirm_password=password123" \
    2>&1)

echo "Response:"
echo "$response"

# Stop the application
kill $APP_PID 2>/dev/null

echo "Test complete."
