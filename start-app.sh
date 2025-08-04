#!/bin/bash

set -e  # Exit on any error

echo "🔧 Complete Registration Fix"
echo "============================"

# Step 1: Start PostgreSQL
echo "1. Setting up PostgreSQL..."
docker rm -f mh-go-postgres 2>/dev/null || true

docker run -d --name mh-go-postgres \
    -e POSTGRES_DB=mh_go \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=password \
    -p 5432:5432 \
    postgres:15-alpine

echo "   Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
    if docker exec mh-go-postgres pg_isready -U postgres > /dev/null 2>&1; then
        echo "   ✅ PostgreSQL is ready"
        break
    fi
    sleep 1
done

# Step 2: Test database connection and setup
echo "2. Testing database setup..."
go run test-db.go

# Step 3: Start the application
echo "3. Starting application..."
echo "   Application will start on http://localhost:8080"
echo "   You can test registration by visiting http://localhost:8080/register"
echo ""
echo "   Press Ctrl+C to stop the application"
echo ""

go run cmd/server/main.go
