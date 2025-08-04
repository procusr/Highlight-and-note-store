# MH-Go Quick Start Guide

## 🚀 Fixed Issues & Ready to Go!

### ✅ Issues Fixed:

1. **Go Version Compatibility**: Updated from Go 1.24.5 to Go 1.21 for Docker compatibility
2. **Docker Build**: Now using Go 1.23-alpine with proper Templ integration
3. **Hot Reloading**: Improved Air configuration with better file watching
4. **Docker Compose**: Added development and production profiles

### 🛠️ Quick Start Options:

#### Option 1: Local Development (Recommended)

```bash
# Copy environment file
cp .env.example .env

# Quick start with script
./dev.sh
```

#### Option 2: Docker Development with Hot Reloading

```bash
# Start development environment
make docker-dev
# OR
docker-compose --profile dev up
```

#### Option 3: Docker Production

```bash
# Start production environment
make docker-prod
# OR
docker-compose --profile prod up --build
```

#### Option 4: Manual Setup

```bash
# Start PostgreSQL only
docker-compose up postgres -d

# Install tools
make install-tools

# Generate templates and start
templ generate
air
```

### 📋 Available Commands:

```bash
# Development
./dev.sh              # Quick start with hot reloading
make dev               # Local development with Air
make docker-dev        # Docker development with hot reloading

# Building
make build             # Build the application
make templ             # Generate Templ templates
go build cmd/server/main.go

# Testing
make test              # Run tests
go test ./...

# Docker
make docker-up         # Start PostgreSQL only
make docker-prod       # Production deployment
make docker-down       # Stop all containers

# Database
make migrate-up        # Run database migrations
make migrate-down      # Rollback migrations
```

### 🌐 Access Points:

- **Application**: http://localhost:8080
- **PostgreSQL**: localhost:5432

### 🗄️ Database Info:

- **Database**: mh_go
- **Username**: postgres
- **Password**: password
- **Connection String**: `postgres://postgres:password@localhost:5432/mh_go?sslmode=disable`

### 🔧 Environment Variables:

Edit `.env` file to customize:

- `PORT`: Server port (default: :8080)
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: JWT signing secret
- `ENVIRONMENT`: development/production

### 📁 Project Features:

- ✅ Authentication (sign-up, sign-in, password reset)
- ✅ Folder-based content organization
- ✅ HTMX dynamic UI with Tailwind CSS
- ✅ Full-text search capability
- ✅ Hot reloading development
- ✅ Docker deployment ready
- ✅ Database migrations
- ✅ Testing framework

### 🔥 Hot Reloading:

The application watches for changes in:

- `.go` files - Go source code
- `.templ` files - Templ templates
- Automatically rebuilds and restarts

### 🎯 Next Steps:

1. Run `./dev.sh` to start development
2. Visit http://localhost:8080
3. Register a new account
4. Start creating folders and content!

Happy coding! 🎉
