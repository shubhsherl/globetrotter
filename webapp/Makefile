.PHONY: setup build start clean test

# Default target
all: setup build

# Setup dependencies
setup:
	@echo "Setting up webapp dependencies..."
	npm install

# Build the application
build:
	@echo "Building webapp..."
	npm run build

# Start the application in development mode
start:
	@echo "Starting webapp in development mode..."
	npm start

# Clean build artifacts
clean:
	@echo "Cleaning webapp..."
	rm -rf build

# Run tests
test:
	@echo "Running tests..."
	npm test 