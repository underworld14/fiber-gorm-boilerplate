# Builder stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags "-s -w -extldflags '-static'" -o /app/bin/server ./cmd/main.go

# Final stage
FROM alpine:3.16

# Add necessary runtime packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/server .

# Create volume for SQLite database persistence
VOLUME ["/app/data"]

# Set environment variables
ENV DB_SOURCE=/app/data/app.db
ENV ENVIRONMENT=production

# Expose port
EXPOSE 3000

# Run the application
CMD ["./server"]
