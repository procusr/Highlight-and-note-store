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

## Project Structure

- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public library code
- `migrations/` - Database migrations
- `templates/` - Templ template files
- `static/` - Static assets (CSS, JS, images)
- `docker/` - Docker configuration files
