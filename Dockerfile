# Use Go 1.23.4 as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app


# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .


# Build the Go application for user service
RUN go build -o users_service ./cmd/http/users/main.go


# Build the Go application for user service
RUN go build -o tasks_service ./cmd/http/tasks/main.go

# Expose the port
EXPOSE 8080 8081
