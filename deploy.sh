#!/bin/bash
# Quick Deployment Script for MH-Go
# This script helps you quickly deploy the containerized application

set -e  # Exit on error

echo "🚀 MH-Go Docker Deployment Script"
echo "=================================="
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    echo "Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    echo "Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

echo -e "${GREEN}✅ Docker and Docker Compose are installed${NC}"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found. Creating from .env.example...${NC}"

    if [ -f .env.example ]; then
        cp .env.example .env
        echo -e "${GREEN}✅ Created .env file${NC}"
        echo ""
        echo -e "${YELLOW}⚠️  IMPORTANT: Please edit .env and update the following:${NC}"
        echo "   - JWT_SECRET (generate with: openssl rand -base64 32)"
        echo "   - Database credentials (if changing defaults)"
        echo "   - SMTP settings (for email functionality)"
        echo "   - OpenAI API key (for AI features)"
        echo ""
        read -p "Press Enter after editing .env file to continue..."
    else
        echo -e "${RED}❌ .env.example not found. Cannot create .env file.${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ .env file exists${NC}"
fi

echo ""
echo "Choose deployment mode:"
echo "1) Development (with hot reload, exposed DB port)"
echo "2) Production (optimized, no exposed DB port)"
read -p "Enter choice [1-2]: " choice

case $choice in
    1)
        echo ""
        echo -e "${GREEN}🔧 Starting Development Environment...${NC}"
        echo ""
        docker-compose up -d --build
        COMPOSE_FILE="docker-compose.yml"
        ;;
    2)
        echo ""
        echo -e "${GREEN}🌐 Starting Production Environment...${NC}"
        echo ""
        docker-compose -f docker-compose.prod.yml up -d --build
        COMPOSE_FILE="docker-compose.prod.yml"
        ;;
    *)
        echo -e "${RED}❌ Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}⏳ Waiting for services to be healthy...${NC}"
echo ""

# Wait for services to be healthy
max_attempts=60
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if docker ps | grep -q "mhgo-app.*healthy" && docker ps | grep -q "mhgo-postgres.*healthy"; then
        echo -e "${GREEN}✅ All services are healthy!${NC}"
        break
    fi

    attempt=$((attempt + 1))
    echo -n "."
    sleep 2

    if [ $attempt -eq $max_attempts ]; then
        echo ""
        echo -e "${RED}❌ Services did not become healthy in time${NC}"
        echo ""
        echo "Check logs with: docker-compose logs"
        exit 1
    fi
done

echo ""
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}🎉 Deployment Successful!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "📊 Service Status:"
docker-compose -f "$COMPOSE_FILE" ps
echo ""
echo "🌐 Application URL: http://localhost:8080"
echo "🏥 Health Check: http://localhost:8080/health"
echo ""
echo "📝 Useful Commands:"
echo "   View logs:       docker-compose logs -f"
echo "   Stop services:   docker-compose down"
echo "   Restart:         docker-compose restart"
echo "   Shell access:    docker exec -it mhgo-app sh"
echo "   DB shell:        docker exec -it mhgo-postgres psql -U user -d mhgo"
echo ""
echo -e "${GREEN}Happy coding! 🚀${NC}"
