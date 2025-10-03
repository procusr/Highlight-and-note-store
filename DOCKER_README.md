# 🐳 Docker Deployment Quick Reference

## Quick Start (3 Steps)

1. **Copy and configure environment variables:**

   ```bash
   cp .env.example .env
   # Edit .env and set your JWT_SECRET
   ```

2. **Run the deployment script:**

   ```bash
   ./deploy.sh
   ```

3. **Access your application:**
   - App: http://localhost:8080
   - Health: http://localhost:8080/health

## Manual Deployment

```bash
# Development mode (with hot reload)
make docker-prod-up

# Production mode
docker-compose -f docker-compose.prod.yml up -d --build

# View logs
make docker-logs

# Stop services
make docker-prod-down
```

## Common Commands

| Command                | Description                       |
| ---------------------- | --------------------------------- |
| `make docker-prod-up`  | Start all services                |
| `make docker-logs`     | View all logs                     |
| `make docker-logs-app` | View app logs only                |
| `make docker-shell`    | Open shell in app container       |
| `make docker-db-shell` | Open PostgreSQL shell             |
| `make docker-ps`       | List running containers           |
| `make docker-restart`  | Restart all services              |
| `make docker-clean`    | Remove all containers and volumes |

## Database Management

```bash
# Access PostgreSQL
make docker-db-shell

# Backup database
make db-backup

# Restore database
make db-restore FILE=backups/backup_2025-10-03.sql

# Reset database (⚠️ deletes all data)
make docker-reset
```

## Architecture

```
┌─────────────────────────────────────┐
│         Docker Network              │
│                                     │
│  ┌──────────────┐  ┌─────────────┐ │
│  │   mhgo-app   │──│  postgres   │ │
│  │  Port: 8080  │  │ Port: 5432  │ │
│  └──────────────┘  └─────────────┘ │
│         │                  │        │
└─────────┼──────────────────┼────────┘
          │                  │
     Host: 8080         Host: 5433
```

## Production Checklist

- [ ] Generate strong JWT_SECRET: `openssl rand -base64 32`
- [ ] Update database credentials in `.env`
- [ ] Configure SMTP settings for email
- [ ] Set up reverse proxy (Nginx/Caddy)
- [ ] Enable HTTPS with Let's Encrypt
- [ ] Configure firewall (only expose 80/443)
- [ ] Set up monitoring and logging
- [ ] Regular database backups

## Troubleshooting

**Container won't start:**

```bash
make docker-logs-app
```

**Database connection issues:**

```bash
docker exec -it mhgo-postgres pg_isready -U user
make docker-logs-db
```

**Port conflicts:**
Edit `docker-compose.yml` and change the host port:

```yaml
ports:
  - "8081:8080" # Changed from 8080
```

## Full Documentation

See [DEPLOYMENT.md](./DEPLOYMENT.md) for comprehensive deployment guide including:

- Cloud deployment strategies
- Security best practices
- Monitoring setup
- Production optimizations
- CI/CD integration

## Support

- 📚 Full docs: [DEPLOYMENT.md](./DEPLOYMENT.md)
- 🐛 Issues: Create a GitHub issue
- 💬 Questions: Check existing discussions
