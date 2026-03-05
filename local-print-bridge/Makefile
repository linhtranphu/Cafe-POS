.PHONY: build run clean test install

# Build binary
build:
	go build -o print-bridge main.go

# Build for different platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o print-bridge-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o print-bridge-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o print-bridge-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o print-bridge-windows-amd64.exe main.go

# Run in development mode
run:
	go run main.go

# Clean build artifacts
clean:
	rm -f print-bridge print-bridge-*

# Run tests
test:
	go test -v ./...

# Install dependencies
install:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Build Docker image
docker-build:
	docker build -t local-print-bridge:latest .

# Run with Docker
docker-run:
	docker run -p 3001:3001 --env-file .env local-print-bridge:latest
