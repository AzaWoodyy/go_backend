# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go module files and download dependencies first to leverage cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application statically (recommended for alpine)
# CGO_ENABLED=0 prevents the build from requiring C libraries in the final image
# -ldflags="-w -s" strips debug information, reducing binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main .

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy .env file - optional, compose handles env vars better usually
# COPY .env .

# Expose the port the app runs on (defined by APP_PORT env var)
# This is documentation; the actual mapping happens in compose.yml
# We don't know the exact port here, but EXPOSE is good practice.
# EXPOSE 8080

# Command to run the executable
# The actual port is determined by the APP_PORT env var passed by compose.yml
CMD ["/app/main"]
