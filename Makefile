.PHONY: build run test clean deps frontend

# Build the binary
build:
	go build -o bin/compareflow cmd/compareflow/main.go

# Run the application
run:
	go run cmd/compareflow/main.go

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build frontend and copy to web/dist
frontend:
	cd ../frontend && npm run build
	rm -rf web/dist
	mkdir -p web/dist
	cp -r ../frontend/dist/* web/dist/

# Build complete application (frontend + backend)
build-all: frontend build

# Clean build artifacts
clean:
	rm -f bin/compareflow
	rm -f compareflow.db
	rm -rf web/dist

# Development mode with air (hot reload)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Create release build
release: frontend
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/compareflow-linux-amd64 cmd/compareflow/main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o bin/compareflow-darwin-amd64 cmd/compareflow/main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o bin/compareflow-darwin-arm64 cmd/compareflow/main.go