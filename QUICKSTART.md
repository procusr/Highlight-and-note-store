# Quick Start Guide

## Running MH-Go Locally

### 1. Easiest Way - Use the Script

```bash
chmod +x start-app.sh
./start-app.sh
```

### 2. Step by Step

```bash
# Start database
docker-compose up -d postgres

# Build and run
go build -o mh-go cmd/server/main.go
./mh-go
```

### 3. Access the Application

Open your browser and go to: **http://localhost:8080**

## First Time Setup

### Create a User Account

1. Go to http://localhost:8080/register
2. Enter your email and password
3. Click "Sign up"

### Start Using the App

1. Create folders to organize your content
2. Add content (text, markdown, links, files) to folders
3. Use the search functionality to find content

## Stopping the Application

- Press `Ctrl+C` in the terminal where the app is running
- Stop the database: `docker-compose down`

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps

# If not running, start it
docker-compose up -d postgres

# Check logs
docker-compose logs postgres
```

### Build Issues

```bash
# Make sure you have Go installed
go version

# Install missing dependencies
go mod tidy

# Regenerate templates if needed
templ generate
```

### Port Already in Use

If port 8080 is busy, set a different port:

```bash
export PORT=3000
./mh-go
```

Then access the app at http://localhost:3000
