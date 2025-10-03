# 🚀 MH-Go Docker Deployment Guide

Complete guide for containerizing and deploying the MH-Go application using Docker.

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Architecture](#architecture)
5. [Configuration](#configuration)
6. [Database Management](#database-management)
7. [Production Deployment](#production-deployment)
8. [Monitoring & Troubleshooting](#monitoring--troubleshooting)
9. [Common Commands](#common-commands)

---

## 🎯 Overview

This deployment uses:

- **Docker multi-stage builds** for optimized Go application images
- **Docker Compose** for orchestrating application and database services
- **PostgreSQL 15** with persistent volumes
- **Automatic database migrations** on application startup
- **Health checks** for both services
- **Non-root user** execution for security

---

## ✅ Prerequisites

- Docker Engine 20.10+ installed
- Docker Compose v2.0+ installed
- 1GB+ available RAM
- 2GB+ available disk space

**Installation:**

```bash
# macOS (Homebrew)
brew install docker docker-compose

# Or download Docker Desktop
https://www.docker.com/products/docker-desktop
```

---

## 🚀 Quick Start

### 1. **Configure Environment Variables**

```bash
# Copy the example environment file
cp .env.example .env

# Edit .env with your configuration
# CRITICAL: Change JWT_SECRET in production!
nano .env
```

**Minimum required changes:**

```bash
JWT_SECRET=$(openssl rand -base64 32)  # Generate a strong secret
```

### 2. **Build and Start Services**

```bash
# Build and start all services
make docker-prod-up

# Or using docker-compose directly
docker-compose up -d --build
```

### 3. **Verify Deployment**

```bash
# Check service status
make docker-ps

# View logs
make docker-logs

# Test the application
curl http://localhost:8080/health
```

**Expected output:**

```json
{ "status": "ok" }
```

### 4. **Access the Application**

Open your browser: **http://localhost:8080**

---

## 🏗️ Architecture

### Container Structure

```
┌─────────────────────────────────────┐
│         Docker Network              │
│        (mhgo-network)               │
│                                     │
│  ┌──────────────┐  ┌─────────────┐ │
│  │   mhgo-app   │  │   postgres  │ │
│  │  (Go App)    │──│ (Database)  │ │
│  │  Port: 8080  │  │ Port: 5432  │ │
│  └──────────────┘  └─────────────┘ │
│         │                  │        │
└─────────┼──────────────────┼────────┘
          │                  │
     Host: 8080         Host: 5433
```

### Multi-Stage Docker Build

**Stage 1: Builder**

- Base: `golang:1.24-alpine`
- Installs templ CLI
- Generates templates
- Compiles Go binary (static, stripped)

**Stage 2: Runtime**

- Base: `alpine:3.19` (~5MB)
- Copies only the binary and assets
- Runs as non-root user
- Final image: ~20-30MB

---

## ⚙️ Configuration

### Environment Variables

All configuration is managed via environment variables in the `.env` file:

| Variable         | Description                    | Default             | Required |
| ---------------- | ------------------------------ | ------------------- | -------- |
| `PORT`           | Application port               | `8080`              | No       |
| `ENVIRONMENT`    | Runtime environment            | `production`        | No       |
| `DATABASE_URL`   | PostgreSQL connection string   | See `.env.example`  | **Yes**  |
| `JWT_SECRET`     | JWT signing secret             | _(must change)_     | **Yes**  |
| `SMTP_HOST`      | Email server host              | `smtp.gmail.com`    | No       |
| `SMTP_PORT`      | Email server port              | `587`               | No       |
| `SMTP_USERNAME`  | Email username                 | -                   | No       |
| `SMTP_PASSWORD`  | Email password                 | -                   | No       |
| `FROM_EMAIL`     | Sender email address           | `noreply@mh-go.com` | No       |
| `OPENAI_API_KEY` | OpenAI API key for AI features | -                   | No       |

### Database Connection Strings

**Local Development (Docker Compose):**

```bash
DATABASE_URL=postgres://user:password@localhost:5433/mhgo?sslmode=disable
```

**Container-to-Container (Production with Docker Compose):**

```bash
DATABASE_URL=postgres://user:password@postgres:5432/mhgo?sslmode=disable
```

**External Database (Cloud/Production):**

```bash
DATABASE_URL=postgres://username:password@your-db-host:5432/dbname?sslmode=require
```

---

## 🗄️ Database Management

### Automatic Migrations

Database migrations run automatically on application startup. Migration files are located in `migrations/`.

### Persistent Data

Database data is stored in a Docker volume: `postgres_data`

```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect mh-go_postgres_data

# Backup database
make db-backup

# Restore database
make db-restore FILE=backups/backup_2025-10-03.sql
```

### Manual Database Access

```bash
# Open PostgreSQL shell
make docker-db-shell

# Or directly
docker exec -it mhgo-postgres psql -U user -d mhgo
```

**Common SQL commands:**

```sql
-- List tables
\dt

-- View users
SELECT * FROM users;

-- Exit
\q
```

### Database Reset (⚠️ Destructive)

```bash
# Delete all data and restart fresh
make docker-reset

# Or manually
docker-compose down -v
docker-compose up -d
```

---

## 🌐 Production Deployment

### 1. **Security Checklist**

- [ ] Generate strong `JWT_SECRET`: `openssl rand -base64 32`
- [ ] Change default database credentials
- [ ] Enable SSL/TLS for database connections (`sslmode=require`)
- [ ] Use secrets management (AWS Secrets Manager, HashiCorp Vault)
- [ ] Configure firewall rules (only expose ports 80/443)
- [ ] Set up reverse proxy (Nginx, Caddy, Traefik)
- [ ] Enable HTTPS with Let's Encrypt
- [ ] Configure log aggregation
- [ ] Set up monitoring and alerts

### 2. **Reverse Proxy Setup (Nginx)**

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. **External Database (AWS RDS, Digital Ocean, etc.)**

Update `docker-compose.yml` to remove the `postgres` service and update the app's `DATABASE_URL`:

```yaml
services:
  app:
    # ... other config ...
    environment:
      DATABASE_URL: postgres://username:password@your-rds-endpoint:5432/dbname?sslmode=require
```

### 4. **Cloud Deployment Options**

#### **AWS (ECS/Fargate)**

```bash
# Push image to ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
docker tag mh-go:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/mh-go:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/mh-go:latest
```

#### **Google Cloud Run**

```bash
# Build and deploy
gcloud builds submit --tag gcr.io/PROJECT-ID/mh-go
gcloud run deploy mh-go --image gcr.io/PROJECT-ID/mh-go --platform managed
```

#### **Digital Ocean App Platform**

```bash
# Connect GitHub repo and configure:
# - Build Command: docker build -t mh-go .
# - Run Command: /app/server
# - Port: 8080
```

#### **Railway / Render / Fly.io**

These platforms auto-detect Dockerfiles. Simply:

1. Connect your Git repository
2. Set environment variables
3. Deploy!

---

## 📊 Monitoring & Troubleshooting

### View Logs

```bash
# All services
make docker-logs

# Application only
make docker-logs-app

# Database only
make docker-logs-db

# Follow logs in real-time
docker-compose logs -f --tail=100
```

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Check container health status
docker ps

# Detailed health check logs
docker inspect mhgo-app | grep -A 10 Health
```

### Common Issues

#### **Container won't start**

```bash
# Check logs
make docker-logs-app

# Common causes:
# - Database not ready (wait for health check)
# - Missing environment variables
# - Port already in use
```

#### **Database connection failed**

```bash
# Verify database is running
docker ps | grep postgres

# Check database logs
make docker-logs-db

# Test connection
docker exec -it mhgo-postgres pg_isready -U user
```

#### **Application errors**

```bash
# Check application logs
make docker-logs-app

# Open shell for debugging
make docker-shell
```

#### **Port conflicts**

```bash
# Change ports in docker-compose.yml
ports:
  - "8081:8080"  # Change host port to 8081
```

### Performance Monitoring

```bash
# View resource usage
docker stats

# Container-specific stats
docker stats mhgo-app mhgo-postgres
```

---

## 🛠️ Common Commands

### Development

```bash
make dev              # Hot reload development
make test             # Run tests
make build            # Build binary
```

### Docker Operations

```bash
make docker-build     # Build images (no cache)
make docker-prod-up   # Build and start services
make docker-prod-down # Stop services
make docker-restart   # Restart services
make docker-ps        # List containers
make docker-logs      # View all logs
make docker-shell     # Open app shell
make docker-db-shell  # Open database shell
make docker-clean     # Remove everything
```

### Database

```bash
make db-backup        # Create backup
make db-restore FILE=backup.sql  # Restore backup
make db-list          # List backups
```

---

## 🔒 Security Best Practices

1. **Never commit `.env` files** to version control
2. **Use strong, unique secrets** for production
3. **Enable SSL/TLS** for all external connections
4. **Run containers as non-root** (already configured)
5. **Keep images updated** regularly
6. **Use read-only file systems** where possible
7. **Implement rate limiting** on API endpoints
8. **Set up log monitoring** and alerts
9. **Regular security audits** and penetration testing
10. **Backup database** regularly

---

## 📝 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)
- [12-Factor App Methodology](https://12factor.net/)

---

## 🆘 Support

For issues and questions:

1. Check the logs: `make docker-logs`
2. Review this documentation
3. Open an issue on GitHub
4. Check existing issues and discussions

---

**Happy Deploying! 🚀**
