#!/bin/bash

# Database backup and restore utility for MH-Go
# Usage: ./db-utils.sh [backup|restore|reset]

DB_CONTAINER="mhgo-postgres"
DB_USER="user"
DB_NAME="mhgo"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

case "$1" in
    "backup")
        echo "Creating database backup..."
        docker exec -t "$DB_CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" --clean --if-exists > "$BACKUP_DIR/mhgo_backup_$TIMESTAMP.sql"
        echo "Backup created: $BACKUP_DIR/mhgo_backup_$TIMESTAMP.sql"
        ;;

    "restore")
        if [ -z "$2" ]; then
            echo "Usage: $0 restore <backup_file>"
            echo "Available backups:"
            ls -la "$BACKUP_DIR"/*.sql 2>/dev/null || echo "No backups found"
            exit 1
        fi

        if [ ! -f "$2" ]; then
            echo "Backup file not found: $2"
            exit 1
        fi

        echo "Restoring database from: $2"
        docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" < "$2"
        echo "Database restored successfully"
        ;;

    "reset")
        echo "Resetting database (this will delete all data)..."
        read -p "Are you sure? Type 'yes' to continue: " confirm
        if [ "$confirm" = "yes" ]; then
            docker exec -it "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "
                DROP SCHEMA public CASCADE;
                CREATE SCHEMA public;
                GRANT ALL ON SCHEMA public TO \"$DB_USER\";
                GRANT ALL ON SCHEMA public TO public;
            "
            echo "Database reset. Run migrations to recreate tables."
        else
            echo "Reset cancelled"
        fi
        ;;

    "list-backups")
        echo "Available backups:"
        ls -la "$BACKUP_DIR"/*.sql 2>/dev/null || echo "No backups found"
        ;;

    *)
        echo "MH-Go Database Utility"
        echo "Usage: $0 [backup|restore|reset|list-backups]"
        echo ""
        echo "Commands:"
        echo "  backup                    - Create a backup of the current database"
        echo "  restore <backup_file>     - Restore database from backup file"
        echo "  reset                     - Reset database (delete all data)"
        echo "  list-backups             - List available backup files"
        echo ""
        echo "Examples:"
        echo "  $0 backup"
        echo "  $0 restore ./backups/mhgo_backup_20250808_100000.sql"
        echo "  $0 list-backups"
        ;;
esac
