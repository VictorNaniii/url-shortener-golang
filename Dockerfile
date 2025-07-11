# syntax=docker/dockerfile:1

# Start from the official Golang image for building
FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/main.go

# Use a minimal base image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/url-shortener .
COPY config/config.go ./config/config.go
EXPOSE 8080
CMD ["./url-shortener"]
