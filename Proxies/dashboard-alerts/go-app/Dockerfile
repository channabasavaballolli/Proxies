# --- Build Stage ---
FROM golang:1.21-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum* ./

# Copy application source code
COPY . .

# Generate go.sum and download dependencies inside the builder container
RUN go mod tidy && go mod download

# Build a statically linked Go binary optimized for container runtimes
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go-app main.go

# --- Run Stage ---
FROM alpine:3.18

# Add non-root user for security compliance (best practice on OpenShift)
RUN adduser -D -u 10001 devopsuser

# Set runtime workspace
WORKDIR /home/devopsuser

# Copy compile binary from builder
COPY --from=builder /go-app ./go-app

# Ensure correct permissions
RUN chown -R 10001:10001 /home/devopsuser

# Run as non-root user
USER 10001

# Expose Go application port
EXPOSE 8080

# Command to execute binary
CMD ["./go-app"]
