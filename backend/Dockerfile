# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the applications
RUN go build -o migrate cmd/migrate/main.go
RUN go build -o core cmd/core/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/migrate .
COPY --from=builder /app/core .

# Copy migrations and documentation
COPY --from=builder /app/database ./database
COPY --from=builder /app/docs ./docs

# Create a directory for the database
RUN mkdir -p /data

# Create entrypoint script with error checking
RUN echo '#!/bin/sh\n\
echo "Starting migrations..."\n\
./migrate up\n\
if [ $? -ne 0 ]; then\n\
  echo "Migration failed!"\n\
  exit 1\n\
fi\n\
echo "Starting application..."\n\
./core' > entrypoint.sh

RUN chmod +x entrypoint.sh

# Expose port (Railway will use PORT env var)
EXPOSE 80

# Default env vars
ENV APP_PORT=80
ENV SQLITE_DATABASE=/data/gelora.db

CMD ["./entrypoint.sh"]
