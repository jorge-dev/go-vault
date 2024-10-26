# Use the latest official Golang image as the base image
FROM golang:1.23.2-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o .build/goVault

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /app/.build/goVault .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./goVault"]
