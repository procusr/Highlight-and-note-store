# MH-Go Web Application

A Go web application using Chi router, Templ templates, HTMX, and PostgreSQL.

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Templ CLI (`go install github.com/a-h/templ/cmd/templ@latest`)

### Running the Application

#### Option 1: Quick Start Script

```bash
./start-app.sh
```

#### Option 2: Manual Steps

```bash
# 1. Start the database
docker-compose up -d postgres

# 2. Generate templates (if modified)
templ generate

# 3. Build and run
go build -o mh-go cmd/server/main.go
./mh-go
```

#### Option 3: Development Mode with Hot Reload

```bash
# Install Air for hot reloading
go install github.com/cosmtrek/air@latest

# Start database
docker-compose up -d postgres

# Run with hot reload
air
```

### Accessing the Application

- **Web Interface**: http://localhost:8080
- **Login Page**: http://localhost:8080/login
- **Register Page**: http://localhost:8080/register

### Default Database Credentials

- **Host**: localhost:5433
- **Database**: mhgo
- **Username**: user
- **Password**: password

## Development

### Project Structure

```
mh-go/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── database/               # Database utilities
│   ├── handlers/               # HTTP request handlers
│   ├── middleware/             # HTTP middleware
│   ├── models/                 # Data models
│   └── services/               # Business logic
├── templates/                  # Templ template files
├── static/                     # Static assets (CSS, JS)
├── migrations/                 # Database migrations
└── docker-compose.yml          # Database setup
```

### Building the Application

#### Standard Build

```bash
go build -o mh-go cmd/server/main.go
```

#### Build for Different Architectures

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o mh-go-linux cmd/server/main.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o mh-go.exe cmd/server/main.go
```

### Template Development

When you modify `.templ` files, regenerate Go code:

```bash
templ generate
```

### Database Operations

```bash
# Start database
docker-compose up -d postgres

# Stop database
docker-compose down

# Reset database (removes all data)
docker-compose down -v
docker-compose up -d postgres
```

### Environment Variables

Create a `.env` file or set these variables:

```bash
export DATABASE_URL="postgres://user:password@localhost:5433/mhgo?sslmode=disable"
export JWT_SECRET="your-super-secret-jwt-key-change-in-production"
export PORT="8080"
export ENVIRONMENT="development"
```

## Features

- User authentication (register/login)
- Folder management system
- Content creation and organization
- HTMX for dynamic interactions
- Responsive design with Tailwind CSS
- PostgreSQL database with migrations

## Tech Stack

- **Backend**: Go with Chi router
- **Templates**: Templ (type-safe HTML templates)
- **Frontend**: HTMX + Tailwind CSS
- **Database**: PostgreSQL
- **Containerization**: Docker
