.PHONY: setup build run clean test migrate

# Variables
DB_PATH=./data/globetrotter.db
GO_FILES=$(shell find . -name "*.go")
MODULE_NAME=github.com/shubhsherl/globetrotter/backend

# Default target
all: setup build run

# Setup dependencies and database
setup:
	@echo "Setting up dependencies..."
	go mod init $(MODULE_NAME) || true
	go get -u github.com/gin-gonic/gin
	go get -u github.com/gin-contrib/cors
	go get -u github.com/mattn/go-sqlite3

# Build the application
build: $(GO_FILES)
	@echo "Building backend..."
	go build -o ./bin/globetrotter

# Run the application
run:
	@echo "Running server..."
	./bin/globetrotter

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f ../bin/globetrotter
	rm -f $(DB_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Initialize database only
init-db:
	@echo "Initializing database..."
	rm -f $(DB_PATH)
	go run cmd/init_db/main.go

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

# Run with hot reload (requires air: https://github.com/cosmtrek/air)
dev:
	@echo "Running with hot reload..."
	air 