# Multi-stage build for Tic-Tac-Go server
# Stage 1: Build
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /build

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary
# -ldflags="-w -s" reduces binary size by stripping debug info
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o tic-tac-go-server \
    ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS support (if needed in future)
# Install wget for health check
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/tic-tac-go-server .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (default 8080, can be overridden via TICTACGO_PORT env var)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${TICTACGO_PORT:-8080}/health || exit 1

# Run the server
# TICTACGO_PORT can be set via environment variable
ENTRYPOINT ["./tic-tac-go-server"]

