# Multi-stage build for TimeSeriesDB
ARG VERSION=dev
ARG GOOS=linux
ARG GOARCH=amd64
ARG PORT=8080
ARG USER_ID=1001
ARG GROUP_ID=1001

# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependenciesgs
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION
ARG GOOS
ARG GOARCH
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -a -installsuffix cgo \
    -o timeseriesdb .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
ARG USER_ID
ARG GROUP_ID
RUN addgroup -g ${GROUP_ID} -S timeseriesdb && \
    adduser -u ${USER_ID} -S timeseriesdb -G timeseriesdb

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/timeseriesdb .

# Copy any additional runtime files
COPY --from=builder /app/env.example .env.example

# Create data directory
RUN mkdir -p /app/data /app/data/backups && \
    chown -R timeseriesdb:timeseriesdb /app

# Switch to non-root user
USER timeseriesdb

# Expose port
ARG PORT
EXPOSE ${PORT}

# Set environment variables
ENV PORT=${PORT}
ENV DATA_FILE=/app/data/data.tsv
ENV DATA_DIR=/app/data
ENV BACKUP_DIR=/app/data/backups

# Health check
HEALTHCHECK --interval=1s --timeout=300s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the application
ENTRYPOINT ["./timeseriesdb"]
