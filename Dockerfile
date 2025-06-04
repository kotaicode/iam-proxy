# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o iam-proxy ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/iam-proxy .

# Run as non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./iam-proxy"] 