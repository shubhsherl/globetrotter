# Backend Tests

This directory contains tests for the Globetrotter backend API.

## Test Structure

The tests are organized by package:

- `api/handlers_test.go`: Tests for API endpoints
- `services/data_service_test.go`: Tests for service layer
- `models/models_test.go`: Tests for data models
- `db/database_test.go`: Tests for database operations

## Running Tests

You can run the tests using the following commands:

### Run all backend tests

```bash
make backend-test
```

Or from the root directory:

```bash
cd backend
make test
```

### Run tests for a specific package

```bash
cd backend
go test ./api -v
go test ./services -v
go test ./models -v
go test ./db -v
```

### Run a specific test

```bash
cd backend
go test ./api -v -run TestHealthCheck
```

## Test Coverage

To generate a test coverage report:

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

This will open a browser window showing the test coverage for each file.

## Adding New Tests

When adding new API endpoints or functionality, please add corresponding tests following the existing patterns.

Test files should be named with the `_test.go` suffix and placed in the same package as the code they're testing. 