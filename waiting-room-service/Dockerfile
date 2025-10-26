# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o waiting-room-service

# --- Run stage ---
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/waiting-room-service .

EXPOSE 8082
CMD ["./waiting-room-service"]
