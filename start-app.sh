#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting MH-Go Application...${NC}"

# Start PostgreSQL with Docker Compose
echo -e "${YELLOW}Starting PostgreSQL database...${NC}"
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
sleep 5

# Check if database is ready
while ! docker exec mhgo-postgres pg_isready -U user -d mhgo; do
  echo -e "${YELLOW}Waiting for database...${NC}"
  sleep 2
done

echo -e "${GREEN}Database is ready!${NC}"

# Run migrations (if you have a migration tool)
echo -e "${YELLOW}Running database migrations...${NC}"
# You can add migration commands here if needed

# Build the application
echo -e "${YELLOW}Building the application...${NC}"
go build -o main .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Build successful!${NC}"

    # Set environment variables
    export DATABASE_URL="postgres://user:password@localhost:5433/mhgo?sslmode=disable"
    export JWT_SECRET="your-super-secret-jwt-key-change-in-production"
    export PORT="8080"

    echo -e "${GREEN}Starting the web server on port 8080...${NC}"
    echo -e "${YELLOW}Visit http://localhost:8080 to access the application${NC}"
    echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"

    # Start the application
    ./main
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
