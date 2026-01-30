# Multi-stage Dockerfile for docker-network-viz
# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -o docker-network-viz \
    .

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/docker-network-viz /app/docker-network-viz

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# The Docker socket must be mounted at runtime:
# docker run -v /var/run/docker.sock:/var/run/docker.sock docker-network-viz

ENTRYPOINT ["/app/docker-network-viz"]
CMD []
