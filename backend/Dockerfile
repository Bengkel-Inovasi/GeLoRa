# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the applications as static binaries
RUN CGO_ENABLED=0 go build -o migrate cmd/migrate/main.go
RUN CGO_ENABLED=0 go build -o seeder cmd/seeder/main.go
RUN CGO_ENABLED=0 go build -o core cmd/core/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/migrate .
COPY --from=builder /app/seeder .
COPY --from=builder /app/core .

# Copy migrations and documentation
COPY --from=builder /app/database ./database
COPY --from=builder /app/docs ./docs

# Create a directory for the database
RUN mkdir -p /data

# Use a direct command string to avoid shell script issues
# This will run migrate up, seeder, and then run core
CMD ["sh", "-c", "./migrate up && ./seeder && APP_PORT=${PORT:-80} ./core"]

# Expose port
EXPOSE 80

# Default env vars
ENV APP_PORT=80
ENV SQLITE_DATABASE=/data/gelora.db
