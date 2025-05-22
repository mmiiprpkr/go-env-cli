.PHONY: build run clean test test-setup docker-up docker-down init-db

# Build variables
APP_NAME=go-env-cli
BUILD_DIR=./build
GO_FILES=$(shell find . -type f -name "*.go")
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "development")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X go-env-cli/cmd.Version=$(VERSION) -X go-env-cli/cmd.GitCommit=$(COMMIT) -X go-env-cli/cmd.BuildDate=$(BUILD_DATE)"

build: $(GO_FILES)
	@echo "Building $(APP_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .

run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME) $(ARGS)

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

test:
	@go test ./...

test-setup: docker-up
	@echo "Setting up test database..."
	@./scripts/setup_test_db.sh

docker-up:
	@echo "Starting PostgreSQL in Docker..."
	@docker-compose up -d

docker-down:
	@echo "Stopping PostgreSQL in Docker..."
	@docker-compose down

init-db:
	@echo "Initializing database..."
	@go run cmd/go-env-cli/init_db.go

# Install the binary to the system
install: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	@cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/

# Help command
help:
	@echo "Available commands:"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application (use ARGS to pass arguments)"
	@echo "  make clean      - Clean build files"
	@echo "  make test       - Run tests"
	@echo "  make test-setup - Set up test database"
	@echo "  make docker-up  - Start PostgreSQL container"
	@echo "  make docker-down- Stop PostgreSQL container"
	@echo "  make init-db    - Initialize database schema"
	@echo "  make install    - Install to /usr/local/bin"
	@echo ""
	@echo "Example usage:"
	@echo "  make run ARGS='list-projects'"
	@echo "  make run ARGS='import .env --project my-app --env development'"
