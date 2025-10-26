# Use an official lightweight Go image
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the Go binary
RUN go build -o ticket-service

# Final small runtime image
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/ticket-service .

# Expose the app port
EXPOSE 8081

# Run the app
CMD ["./ticket-service"]
