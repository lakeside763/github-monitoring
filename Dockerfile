FROM golang:1.23

WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o main ./cmd/app/main.go

# Expose port
EXPOSE 5200

# Run the binary
CMD ["./main"]