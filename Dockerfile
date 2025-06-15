# Build stage
FROM golang:1.23-alpine AS builder

# Install git for version info
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG COMMIT_HASH=unknown
ARG BUILD_DATE=unknown

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w \
    -X github.com/gigiozzz/kubectl-hello-world/cmd.Version=${VERSION} \
    -X github.com/gigiozzz/kubectl-hello-world/cmd.CommitHash=${COMMIT_HASH} \
    -X github.com/gigiozzz/kubectl-hello-world/cmd.BuildDate=${BUILD_DATE}" \
    -o kubectl-hello-world .

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 kubectl && \
    adduser -D -u 1001 -G kubectl kubectl

# Set working directory
WORKDIR /home/kubectl

# Copy binary from builder stage
COPY --from=builder /app/kubectl-hello-world /usr/local/bin/kubectl-hello-world

# Make sure binary is executable
RUN chmod +x /usr/local/bin/kubectl-hello-world

# Switch to non-root user
USER kubectl

# Create .kube directory for kubeconfig mounting
RUN mkdir -p /home/kubectl/.kube

# Add metadata
LABEL org.opencontainers.image.title="kubectl-hello-world"
LABEL org.opencontainers.image.description="A kubectl plugin demonstrating best practices"
LABEL org.opencontainers.image.url="https://github.com/gigiozzz/kubectl-hello-world"
LABEL org.opencontainers.image.source="https://github.com/gigiozzz/kubectl-hello-world"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_DATE}"
LABEL org.opencontainers.image.revision="${COMMIT_HASH}"
LABEL org.opencontainers.image.licenses="GPL 3.0"

# Default command
ENTRYPOINT ["kubectl-hello-world"]
CMD ["--help"]