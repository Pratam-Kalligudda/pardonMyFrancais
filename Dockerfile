# Use an official Golang runtime as a parent image
FROM golang:1.21-alpine

# Set the working directory in the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

RUN apk update && \
    apk add --no-cache ffmpeg
    
# Install any dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./main"]
