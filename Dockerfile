# Stage 1: Build the React frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/webapp

# Copy package.json and package-lock.json
COPY webapp/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the frontend code
COPY webapp/ ./

# Build the frontend
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Install SQLite dependencies for the build
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go.mod and go.sum
COPY backend/go.mod ./
COPY backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the backend source code
COPY backend/ ./

# Build the backend
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/globetrotter github.com/shubhsherl/globetrotter/backend

# Stage 3: Final image
FROM alpine:3.18

# Install necessary packages
RUN apk add --no-cache ca-certificates tzdata sqlite curl

# Set working directory
WORKDIR /app

# Create data directory
RUN mkdir -p /app/data

# Copy the built frontend from the frontend-builder stage
COPY --from=frontend-builder /app/webapp/build /app/webapp/build

# Copy the built backend binary from the backend-builder stage
COPY --from=backend-builder /app/bin/globetrotter /app/globetrotter

# Copy migrations if needed
COPY backend/migrations /app/backend/migrations

# Set environment variables
ENV PORT=8080
ENV DB_PATH=/app/data/globetrotter.db
ENV GIN_MODE=release

# Expose the port
EXPOSE 8080

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/globetrotter"] 