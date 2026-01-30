#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting MH-Go Application...${NC}"

# Function to check if command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

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

# Generate templ templates (required for templates.* components)
echo -e "${YELLOW}Generating Templ templates...${NC}"
if command_exists templ; then
  templ generate
else
  echo -e "${YELLOW}Templ not found. Installing...${NC}"
  go install github.com/a-h/templ/cmd/templ@latest
  if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to install templ. Make sure Go is installed and GOPATH/bin is on your PATH.${NC}"
    echo -e "${YELLOW}Tip: run ./setup-dev.sh to install required dev tools.${NC}"
    exit 1
  fi
  templ generate
fi

if [ $? -ne 0 ]; then
  echo -e "${RED}Failed to generate templ templates. Fix template errors and try again.${NC}"
  exit 1
fi

# Build the application
echo -e "${YELLOW}Building the application...${NC}"
mkdir -p bin
go build -o bin/server cmd/server/main.go

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
    ./bin/server
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi
