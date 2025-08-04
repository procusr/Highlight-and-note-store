# MH-Go: Content Management System

A modern, minimalistic content management system built with Go, featuring folder-based organization and real-time updates.

## Features

- 🚀 **Fast & Lightweight**: Built with Go and Chi router
- 🎨 **Modern UI**: HTMX + Tailwind CSS for dynamic, responsive interface
- 📁 **Folder Organization**: Organize content in customizable folders
- 🔐 **Authentication**: Complete auth system with sign-up, sign-in, password reset
- 🔍 **Full-text Search**: PostgreSQL-powered search with tsvector
- ⚡ **Hot Reloading**: Development with Air for instant reloads
- 🐳 **Docker Ready**: Complete containerization setup
- 🗄️ **Database Migrations**: Automated schema management
- ✅ **Testing**: Comprehensive test suite setup

## Tech Stack

- **Backend**: Go 1.21+, Chi Router, PostgreSQL
- **Frontend**: HTMX, Tailwind CSS
- **Database**: PostgreSQL with migrations
- **Development**: Air (hot reloading), Docker
- **Authentication**: JWT tokens with secure password hashing

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional, for convenience commands)

### Development Setup

1. **Quick start with script**:

   ```bash
   git clone <repository>
   cd mH-go
   cp .env.example .env
   ./dev.sh
   ```

2. **Manual setup**:

   ```bash
   # Start PostgreSQL
   docker-compose up postgres -d

   # Install development tools
   make install-tools

   # Generate templates and start with hot reloading
   templ generate
   air
   ```

3. **Docker development with hot reloading**:
   ```bash
   make docker-dev
   ```

The application will be available at `http://localhost:8080`

### Manual Setup

1. **Start PostgreSQL**:

   ```bash
   docker-compose up postgres -d
   ```

2. **Run migrations**:

   ```bash
   make migrate-up
   ```

3. **Install dependencies**:

   ```bash
   go mod download
   ```

4. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

## Usage

### Authentication

1. Navigate to `http://localhost:8080`
2. Register a new account or sign in
3. Use the forgot password feature if needed

### Content Management

1. **Folders**: Create and organize your content in folders
2. **Content**: Add, edit, and manage your content items
3. **Search**: Use full-text search to find content quickly
4. **Favorites**: Mark important content as favorites

## API Endpoints

### Authentication

- `GET /login` - Login page
- `POST /login` - Authenticate user
- `GET /register` - Registration page
- `POST /register` - Create new user
- `POST /logout` - Sign out user
- `GET /forgot-password` - Forgot password page
- `POST /forgot-password` - Request password reset
- `GET /reset-password` - Reset password page
- `POST /reset-password` - Reset password

### Folders (Protected)

- `GET /folders` - List user folders
- `POST /folders` - Create new folder
- `GET /folders/{id}` - Get folder details
- `PUT /folders/{id}` - Update folder
- `DELETE /folders/{id}` - Delete folder

### Content (Protected)

- `GET /content` - List content (optionally filtered by folder)
- `POST /content` - Create new content
- `GET /content/{id}` - Get content details
- `PUT /content/{id}` - Update content
- `DELETE /content/{id}` - Delete content
- `POST /content/{id}/archive` - Archive content
- `POST /content/{id}/favorite` - Toggle favorite status

## Database Schema

The application uses the provided schema with:

- `auth.users` - User authentication
- `public.user_folders` - Content organization
- `public.user_content` - Content storage with full-text search

## Development

### Available Commands

```bash
make build          # Build the application
make run            # Run the application
make dev            # Run with hot reloading (local)
make test           # Run tests
make docker-dev     # Development with Docker + hot reloading
make docker-prod    # Production with Docker
make docker-up      # Start PostgreSQL only
make docker-down    # Stop Docker containers
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations
make templ          # Generate Templ templates
make fmt            # Format code
```

### Project Structure

```
mH-go/
├── cmd/server/         # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database connection & migrations
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Data models
│   └── services/       # Business logic
├── migrations/         # Database migrations
├── static/            # Static files
├── templates/         # Template files
├── tests/             # Test files
├── docker/            # Docker configuration
├── .air.toml          # Hot reload configuration
├── Dockerfile         # Container definition
├── docker-compose.yml # Development environment
└── Makefile          # Development commands
```

### Adding New Features

1. **Models**: Add new data structures in `internal/models/`
2. **Services**: Implement business logic in `internal/services/`
3. **Handlers**: Create HTTP handlers in `internal/handlers/`
4. **Migrations**: Add database changes in `migrations/`
5. **Tests**: Add tests in `tests/`

## Deployment

### Docker Deployment

```bash
# Build and start with Docker
docker-compose up --build -d

# Or build for production
docker build -t mh-go .
docker run -p 8080:8080 \
  -e DATABASE_URL="your-postgres-url" \
  -e JWT_SECRET="your-secret" \
  mh-go
```

### Environment Variables

- `PORT` - Server port (default: `:8080`)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `ENVIRONMENT` - Environment mode (`development`/`production`)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details
