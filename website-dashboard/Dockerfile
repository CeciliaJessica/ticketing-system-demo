FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o website-dashboard .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/website-dashboard .
COPY static ./static
EXPOSE 8083
CMD ["./website-dashboard"]
