# Use the latest stable Go image with Alpine
FROM golang:1.23-alpine

# Install necessary dependencies
# RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy only module files first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install CompileDaemon for live reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

# Copy the rest of the app
COPY . .

# Run CompileDaemon to auto-restart on code changes
CMD ["CompileDaemon", "--build", "go build -o main .", "--command", "./main"]

