FROM golang:1.25-alpine  AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o dashboard-service

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/dashboard-service .
EXPOSE 8083
CMD ["./dashboard-service"]
