BINARY_NAME=http-server

.PHONY: all build run test clean help migrate

all: build

check:
	@go fmt ./...
	@golangci-lint run -v ./...
	@govulncheck ./...

build:
	@echo "Generating swagger docs"
	@swag init --generalInfo cmd/api/main.go --output docs
	@echo "Building the application..."
	@mkdir -p build
	@go build -o build/$(BINARY_NAME) cmd/api/main.go
	@go build -o build/migrate cmd/migrate/main.go

run: build
	@echo "Running the application..."
	@./build/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -f build/$(BINARY_NAME)

# Migrations
migrate:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go

# Docker
docker: build
	@docker build -t tools/simple-http-server .

run-docker:
	@docker run -it --name simple-http-server -p 8080:8080 tools/simple-http-server

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  check       Run checks for the code base"
	@echo "  build       Build the application"
	@echo "  run         Run the application"
	@echo "  test        Run the tests"
	@echo "  clean       Clean the build artifacts"
	@echo "  migrate     Run database migrations"
	@echo "  docker      Build docker image"
	@echo "  run-docker  Run docker image" 
	@echo "  help     Display this help message"

.DEFAULT_GOAL := help
