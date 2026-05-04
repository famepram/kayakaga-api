# =========================
# Build stage
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Install swag CLI (untuk generate docs)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate swagger docs (INI YANG PENTING)
RUN swag init

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o kayakaga-api .

# =========================
# Runtime stage
# =========================
FROM alpine:latest

WORKDIR /app

# Install CA certificates (HTTPS)
RUN apk --no-cache add ca-certificates

# Copy binary
COPY --from=builder /app/kayakaga-api .

# Expose port
EXPOSE 8080

# Run app
CMD ["./kayakaga-api"]