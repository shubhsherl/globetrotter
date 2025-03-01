.PHONY: setup build run clean test dev webapp-build webapp-start backend-build backend-run all

# Default target
all: setup build run

# Setup both webapp and backend
setup: webapp-setup backend-setup

# Build both webapp and backend
build: webapp-build backend-build

# Run both webapp and backend
run: backend-run webapp-start

# Clean both webapp and backend
clean: webapp-clean backend-clean

# Setup webapp
webapp-setup:
	@echo "Setting up webapp dependencies..."
	cd webapp && npm install

# Build webapp
webapp-build:
	@echo "Building webapp..."
	cd webapp && npm run build

# Start webapp in development mode
webapp-start:
	@echo "Starting webapp in development mode..."
	cd webapp && npm start

# Clean webapp
webapp-clean:
	@echo "Cleaning webapp..."
	rm -rf webapp/build

# Setup backend
backend-setup:
	@echo "Setting up backend..."
	cd backend && make setup

# Build backend
backend-build:
	@echo "Building backend..."
	cd backend && make build

# Run backend
backend-run:
	@echo "Running backend..."
	cd backend && make run

# Clean backend
backend-clean:
	@echo "Cleaning backend..."
	cd backend && make clean

# Run tests for both webapp and backend
test:
	@echo "Running webapp tests..."
	cd webapp && npm test
	@echo "Running backend tests..."
	cd backend && make test

# Run both in development mode
dev:
	@echo "Running in development mode..."
	cd backend && make dev & cd webapp && npm start
