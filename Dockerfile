# Use Go 1.23.4 as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app



# Install required dependencies for CGO
RUN apk add --no-cache gcc musl-dev

# Set CGO enabled
ENV CGO_ENABLED=1

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Set the environment variable
ENV DB_URL="root:root@tcp(localhost:3306)/task_manager?parseTime=true" 
ENV JWT_SEC="ThisIsSec" 
ENV JWT_REF_SEC="ThisIsSuperSec"
ENV JWT_EXP_DURATION="24h"
ENV REF_EXP_DURATION="48h"

# Build the Go application
RUN go build -o main .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"]
