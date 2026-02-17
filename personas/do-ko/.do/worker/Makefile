.PHONY: build run test clean deps

VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT)"

# Build the worker binary
build:
	go build $(LDFLAGS) -o bin/do-worker ./cmd/worker

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/do-worker-darwin-amd64 ./cmd/worker
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/do-worker-darwin-arm64 ./cmd/worker
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/do-worker-linux-amd64 ./cmd/worker
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/do-worker-linux-arm64 ./cmd/worker

# Run the worker
run:
	go run ./cmd/worker

# Run with hot reload (requires air)
dev:
	air

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Check for linting issues
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Generate mocks (requires mockgen)
mocks:
	mockgen -source=internal/db/adapter.go -destination=internal/db/mock_adapter.go -package=db

# Health check
health:
	curl -s http://localhost:3778/health | jq .
