# src/docker/Dockerfile

# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download the dependencies
COPY ../go.mod ../go.sum ./
RUN go mod download

# Copy the application source code
COPY ../ .

# Build the Go application
RUN env GOOS=linux GOARCH=amd64 go build -o backend ./cmd

# Expose the port the app will run on
EXPOSE 8080

# Run the Go application
CMD ["./backend"]