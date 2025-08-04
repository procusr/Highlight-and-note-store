# Copilot Instructions

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

This is a Go web application using:

- Chi router for HTTP routing
- Templ for HTML template generation
- HTMX for dynamic frontend interactions
- Tailwind CSS for styling
- PostgreSQL for database with migrations
- Air for hot reloading during development
- Docker for containerization
- Built-in authentication system

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
