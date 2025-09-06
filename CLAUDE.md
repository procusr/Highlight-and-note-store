# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running

- `make build` - Build the application to `bin/server`
- `make run` - Build and run the application
- `make dev` - Run with hot reload using Air (preferred for development)
- `air` - Direct hot reload command (watches .templ, .go files and auto-generates templates)

### Testing

- `make test` - Run all tests with `go test ./...`
- `go test ./...` - Direct test command

### Database Management

- `make docker-up` - Start PostgreSQL database (localhost:5433)
- `make docker-down` - Stop database services
- `make docker-reset` - Reset database and volumes (deletes all data)
- `make db-backup` - Create database backup
- `make db-restore FILE=<file>` - Restore from backup
- `make db-reset` - Reset database schema and data

### Template Development

- `templ generate` - Generate Go code from .templ files (auto-run by Air during development)

## Architecture Overview

This is a Go web application using Clean Architecture principles with the following layers:

### Core Structure

```
cmd/server/main.go          # Application entry point with dependency injection
internal/
├── config/                 # Configuration management (env vars, JWT secrets)
├── database/               # Database connection, migrations, utilities
├── models/                 # Data models and domain entities
├── services/               # Business logic layer (auth, user, folder, content)
├── handlers/               # HTTP handlers (presentation layer)
└── middleware/             # HTTP middleware (auth, logging)
templates/                  # Templ HTML templates (compiled to Go)
migrations/                 # SQL database migrations
static/                     # CSS, JS, images
```

### Key Services

- **AuthService**: JWT-based authentication, password hashing
- **UserService**: User management operations
- **FolderService**: Folder hierarchy and management
- **ContentService**: Content creation and organization

### Technology Stack

- **Router**: Chi v5 with middleware
- **Templates**: Templ (type-safe HTML templates compiled to Go)
- **Frontend**: HTMX + Alpine.js + Tailwind CSS
- **Database**: PostgreSQL with golang-migrate
- **Authentication**: JWT with bcrypt password hashing
- **Hot Reload**: Air with automatic template generation

### Development Flow

1. Modify .templ files for UI changes
2. Air automatically runs `templ generate` and rebuilds
3. Database migrations run automatically on startup
4. Use dependency injection pattern throughout

### Environment Variables

Create `.env` or set:

```
DATABASE_URL=postgres://user:password@localhost:5433/mhgo?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
ENVIRONMENT=development
```

### Database Connection

- Default: localhost:5433 (Docker)
- Database: mhgo
- User: user / password: password
- Automatic migrations on startup

## Code Conventions

Follow the patterns established in the codebase:

- Use dependency injection for services
- Implement proper error handling with wrapped errors
- Follow Clean Architecture separation (handlers → services → models)
- Use interfaces for testability
- Structure handlers with consistent patterns for HTMX responses
- Use Templ components for reusable UI elements

---

- **Ensure all responses are:**
  - Accurate and factually correct
  - Thoroughly researched and comprehensive
  - Anticipatory of follow-up needs or potential issues

## 🔍 Your Core Identity

You are a technical expert with advanced knowledge in:

- **Frontend Technologies**: Tailwind CSS,Alpinejs,Htmx,UX/UI,
- **Backend Technologies**: Node.js, Golang, microservices architecture
- **Development Practices**: Design patterns, clean architecture, modular development
  Your responses must be thoughtful, nuanced, technically accurate, and demonstrate exceptional reasoning ability.

## 📋 General Coding Principles

**Always follow these core principles:**

1. Provide concise, technical responses with practical examples
2. Favor functional, declarative programming approaches
3. Emphasize modular design and eliminate duplication
4. Enforce separation of concerns through Clean Architecture
5. Ensure code is idiomatic, testable, and follows modern best practices
6. Include complete, working examples
7. Explain key design decisions and trade-offs
8. Address edge cases and error handling

## 🔧 Technology-Specific Guidelines

## 📝 Response Format

Structure all responses with appropriate headings, code blocks, and formatting to maximize clarity and readability. Use numbered lists for sequential steps and bullet points for related concepts.

This is a Go web application using:

- Chi router for HTTP routing
- Templ for HTML template generation
- HTMX for dynamic frontend interactions
- Tailwind CSS for styling
- PostgreSQL for database with migrations
- Air for hot reloading during development
- Docker for containerization
- Built-in authentication system
- Alpinejs for frontend interactivity

## Code Style Guidelines

- Follow standard Go conventions and idioms
- Use dependency injection for better testability
- Implement proper error handling with custom error types
- Structure code using clean architecture principles
- Use interfaces for external dependencies
- Write comprehensive tests for all business logic

* **Concurrency Management**
    - Implement safe goroutine practices with proper cancellation
    - Use context propagation to prevent leaks and deadlocks
    - Guard shared state with channels or sync primitives
* **Architecture & Design**
    - Follow interface-driven development with explicit dependency injection
    - Prefer composition over inheritance
    - Design with clear boundaries and dependency inversion
    - Avoid global state; use constructor functions for dependencies
* **Error Handling & Resources**
    - Check and handle errors explicitly with wrapped errors for traceability
    - Use defer for resource cleanup to prevent leaks
    - Implement robust logging and observability
* **Dependencies**
    - Prefer standard library when feasible
    - Use Go modules with version-locked dependencies

## Project Structure

- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public library code
- `migrations/` - Database migrations
- `templates/` - Templ template files
- `static/` - Static assets (CSS, JS, images)
- `docker/` - Docker configuration files
