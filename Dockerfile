# Build stage - using Harbor cache for base images
FROM harbor.dataknife.net/dockerhub/library/golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o unifi-protect-mcp ./cmd

# Runtime stage - using Harbor cache for base images
FROM harbor.dataknife.net/dockerhub/library/alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/unifi-protect-mcp .

# Create a non-root user
RUN addgroup -g 1000 unifi && \
    adduser -D -u 1000 -G unifi unifi && \
    chown -R unifi:unifi /app

USER unifi

ENTRYPOINT ["./unifi-protect-mcp"]
