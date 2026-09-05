# ============================================
# brain-cli Dockerfile (Multi-stage Build)
# ============================================
# Stage 1: Builder
# Stage 2: Runtime
# ============================================

# ============================================
# STAGE 1: Builder
# ============================================
FROM golang:1.23-alpine AS builder

# Metadata
LABEL maintainer="Nero <your-email@example.com>"
LABEL description="Brain CLI - Asistente Personal TUI con IA"
LABEL version="2.0.0"

# Install build dependencies (needed for CGO and SQLite)
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    git \
    make

# Set working directory
WORKDIR /build

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO enabled (required for SQLite)
# -ldflags: Strip debug info and set version
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s -X main.Version=2.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o brain-cli \
    ./cmd/brain-cli

# Verify binary
RUN chmod +x brain-cli && \
    ./brain-cli --version || true

# ============================================
# STAGE 2: Runtime
# ============================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    bash \
    curl \
    git \
    sqlite \
    sqlite-libs \
    ca-certificates \
    tzdata \
    python3 \
    py3-pip \
    openvpn \
    nmap-ncat \
    docker-cli \
    && rm -rf /var/cache/apk/*

# Create non-root user for security
RUN addgroup -g 1000 brain && \
    adduser -D -u 1000 -G brain -s /bin/bash brain

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=brain:brain /build/brain-cli .

# Copy embedded resources
COPY --from=builder --chown=brain:brain /build/embed ./embed
COPY --from=builder --chown=brain:brain /build/configs ./configs

# Create directories for volumes
RUN mkdir -p /app/data /app/scripts /app/logs && \
    chown -R brain:brain /app

# Switch to non-root user
USER brain

# Environment variables
ENV DB_PATH=/app/data/brain.db \
    CONFIG_PATH=/app/configs/config.yaml \
    LOG_LEVEL=info \
    TERM=xterm-256color \
    NO_COLOR=false

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD [ "/app/brain-cli", "health" ] || exit 1

# Expose port for OmniRoute (if needed)
EXPOSE 20128

# Volume mount points
VOLUME ["/app/data", "/app/scripts", "/app/logs"]

# Entrypoint
ENTRYPOINT ["/app/brain-cli"]

# Default command
CMD ["--help"]
