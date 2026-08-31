# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency files and source code
COPY go.mod ./
COPY main.go ./

# Download dependencies and generate go.sum
RUN go mod tidy

# Build the Go application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o sandbox-backend .

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and pg_isready equivalent
RUN apk --no-cache add ca-certificates postgresql-client

# Copy binary from builder
COPY --from=builder /app/sandbox-backend .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./sandbox-backend"]
