# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and other dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o wallet-api ./cmd/main.go

# Stage 2: Run the Go binary
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/wallet-api .

# Copy the .env file if needed
COPY .env .env

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./wallet-api"]