# Use the official Golang image as a base for building
FROM golang:1.22-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for static binaries, useful for scratch images
# -o specifies the output file name
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-management ./main.go

# Use a minimal image for the final stage
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/notification-management .

# Expose the port the application listens on (assuming 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./notification-management"]
