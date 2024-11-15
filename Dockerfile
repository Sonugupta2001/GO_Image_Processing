# Use Go base image
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Expose application port
EXPOSE 8080

# Default command
CMD ["./main"]
