# Step 1: Build the Go app
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app (this time without specifying GOOS and GOARCH explicitly)
RUN go build -o main .

# Step 2: Create a smaller image to run the app
FROM alpine:latest  

# Install required dependencies (e.g., for database or S3 client)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /app/main .

# Expose the port that the app will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
