# ==================================
# Production Dockerfile for issue2md
# Multi-stage build with security best practices
# ==================================

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies (will be cached if go.mod/go.sum don't change)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version and build info
ARG VERSION=dev
ARG BUILD_TIME
ARG COMMIT_SHA

# Build CLI application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -ldflags="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commitSHA=${COMMIT_SHA}" \
    -o issue2md ./cmd/issue2md

# Build Web application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -ldflags="-X main.webVersion=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.commitSHA=${COMMIT_SHA}" \
    -o issue2mdweb ./cmd/issue2mdweb

# Final stage - minimal runtime image
FROM alpine:3.18

# Install runtime dependencies only
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

# Create non-root user with proper UID/GID
RUN addgroup -g 1001 -S issue2md && \
    adduser -u 1001 -S issue2md -G issue2md -h /app -s /sbin/nologin

# Set working directory
WORKDIR /app

# Copy CA certificates from builder (if needed)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy compiled binaries from builder stage
COPY --from=builder /build/issue2md .
COPY --from=builder /build/issue2mdweb .

# Create directory for potential output files
RUN mkdir -p /app/output && \
    chown -R issue2md:issue2md /app

# Switch to non-root user
USER issue2md

# Expose port for web service
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV GITHUB_TOKEN=""

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Multi-purpose entrypoint that can run either CLI or web service
ENTRYPOINT ["/app/entrypoint.sh"]

# Default to web service
CMD ["web"]

# Optional: Create entrypoint script to handle both CLI and web modes
# This will be copied in a separate instruction to ensure proper permissions
COPY --chown=issue2md:issue2md docker-entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh