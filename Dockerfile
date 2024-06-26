# Use the official Golang base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application executable
RUN go build -o vc .

# Expose the port used by the HTTP Server
EXPOSE 8080

# Default command (can be overridden by Docker Compose)
CMD ["./vc"]

