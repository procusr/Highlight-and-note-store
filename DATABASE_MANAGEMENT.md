# Database Management & Deployment Guide

## Avoiding Data Loss in the Future

### 1. **Local Development Best Practices**

#### Use Docker Volumes (Already Set Up)

Your `docker-compose.yml` uses named volumes which persist data even when containers are recreated:

```yaml
volumes:
  - postgres_data:/var/lib/postgresql/data
```

#### Regular Backups

```bash
# Create a backup before major changes
make db-backup

# Or use the utility directly
./db-utils.sh backup
```

#### Check Volume Status

```bash
# List Docker volumes
docker volume ls

# Inspect the postgres volume
docker volume inspect mh-go_postgres_data

# Check volume location
docker volume inspect mh-go_postgres_data | grep Mountpoint
```

### 2. **Before Switching Docker Environments**

```bash
# 1. Create a backup
./db-utils.sh backup

# 2. Note your current volume
docker volume inspect mh-go_postgres_data

# 3. After switching Docker, restore data
./db-utils.sh restore ./backups/mhgo_backup_YYYYMMDD_HHMMSS.sql
```

### 3. **Production Deployment Options**

#### Option A: Managed Database (Recommended)

- **PostgreSQL providers**: AWS RDS, Google Cloud SQL, DigitalOcean Managed Databases
- **Benefits**: Automated backups, high availability, managed updates
- **Configuration**: Update `DATABASE_URL` environment variable

#### Option B: Self-Hosted with Backup Strategy

```yaml
# docker-compose.prod.yml
services:
  postgres:
    image: postgres:15
    volumes:
      - ./data/postgres:/var/lib/postgresql/data # Host directory
      - ./backups:/backups # Backup directory
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
```

#### Option C: Cloud Deployment with Volumes

```bash
# For cloud providers like DigitalOcean, AWS ECS, etc.
# Use persistent volumes or managed database services
```

### 4. **Migration Strategy for Production**

#### Pre-Deployment Checklist

1. **Test migrations locally**:

   ```bash
   make db-backup  # Backup current state
   make db-reset   # Reset to clean state
   make dev        # Run app, migrations should apply
   ```

2. **Version your database schema**:

   - Keep migrations in version control
   - Test migration rollbacks with `.down.sql` files

3. **Production deployment**:
   ```bash
   # On production server
   DATABASE_URL="postgres://user:pass@prod-db:5432/mhgo" ./bin/server
   ```

### 5. **Environment Variables for Different Stages**

Create environment-specific config files:

#### `.env.local` (Development)

```bash
DATABASE_URL=postgres://user:password@localhost:5433/mhgo?sslmode=disable
JWT_SECRET=dev-secret-key
ENVIRONMENT=development
```

#### `.env.production` (Production)

```bash
DATABASE_URL=postgres://produser:securepass@prod-db.example.com:5432/mhgo?sslmode=require
JWT_SECRET=super-secure-production-key
ENVIRONMENT=production
```

### 6. **Automated Backup Strategy**

#### Daily Backups (Cron Job)

```bash
# Add to crontab: crontab -e
0 2 * * * cd /path/to/mh-go && ./db-utils.sh backup
```

#### Cloud Storage Backup

```bash
# After creating backup, upload to cloud
aws s3 cp ./backups/mhgo_backup_$(date +%Y%m%d_%H%M%S).sql s3://your-backup-bucket/
```

### 7. **Docker Volume Management**

#### Export Volume Data

```bash
# Create a backup of the entire volume
docker run --rm -v mh-go_postgres_data:/data -v $(pwd):/backup ubuntu tar czf /backup/postgres_volume_backup.tar.gz -C /data .
```

#### Import Volume Data

```bash
# Restore volume from backup
docker run --rm -v mh-go_postgres_data:/data -v $(pwd):/backup ubuntu tar xzf /backup/postgres_volume_backup.tar.gz -C /data
```

### 8. **Quick Recovery Commands**

```bash
# Emergency: Start fresh and restore from backup
make docker-reset
make docker-up
sleep 10  # Wait for PostgreSQL to start
./db-utils.sh restore ./backups/your_backup.sql

# Check if data is restored
docker exec -it mhgo-postgres psql -U user -d mhgo -c "SELECT COUNT(*) FROM users;"
```

### 9. **Monitoring & Alerts**

Consider adding health checks and monitoring:

- Database connectivity checks
- Automated backup verification
- Disk space monitoring for backup storage
- Migration status monitoring

This guide should help you avoid data loss scenarios in the future!
