# -------- Build stage --------
FROM golang:1.23.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata build-base

WORKDIR /app

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /main ./cmd/server/main.go

# -------- Runtime stage --------
FROM alpine:3.20

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Non-root user (optional, but recommended)
RUN addgroup -g 1001 -S appgroup && \
  adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appgroup /main ./main

USER appuser

EXPOSE 8080

# Run your service
CMD ["./main"]
