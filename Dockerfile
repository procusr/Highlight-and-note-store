# ============================================
# Stage 1: Build Stage
# ============================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Install templ CLI for template generation
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set working directory
WORKDIR /app

# Copy dependency files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ templates
RUN templ generate

# Build the Go application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to strip debug info and reduce size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/server \
    ./cmd/server/main.go

# ============================================
# Stage 2: Runtime Stage
# ============================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy CA certificates and timezone data from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from builder
COPY --from=builder --chown=appuser:appuser /app/bin/server /app/server

# Copy static assets and migrations
COPY --chown=appuser:appuser static ./static
COPY --chown=appuser:appuser migrations ./migrations

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/server"]
