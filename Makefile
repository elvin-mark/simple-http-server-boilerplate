BINARY_NAME=http-server

.PHONY: all build run test clean help

all: build

build:
	@echo "Building the application..."
	@mkdir -p build
	@go build -o build/$(BINARY_NAME) main.go

run: build
	@echo "Running the application..."
	@./build/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -f build/$(BINARY_NAME)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build    Build the application"
	@echo "  run      Run the application"
	@echo "  test     Run the tests"
	@echo "  clean    Clean the build artifacts"
	@echo "  help     Display this help message"

.DEFAULT_GOAL := help
