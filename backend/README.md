# Globetrotter Challenge Backend

This is the backend API for the Globetrotter Challenge game, built with Go (Golang), the Gin web framework, and SQLite database.

## Features

- RESTful API for game data and user management
- Random destination selection with multiple-choice options
- User score tracking
- Challenge sharing functionality
- SQLite database for persistent storage

## Prerequisites

- Go 1.16 or higher
- Git
- SQLite3 (included as a Go dependency)

## Installation

1. Clone the repository (if you haven't already):
   ```bash
   git clone https://github.com/yourusername/globetrotter.git
   cd globetrotter
   ```

2. Set up the project using the Makefile:
   ```bash
   make setup
   ```

3. Build the application:
   ```bash
   make build
   ```

4. Run the server:
   ```bash
   make run
   ```

   The server will start on port 8080 by default.

## Database Management

The application uses SQLite for data storage. The database will be automatically initialized when you run the server for the first time.

To reinitialize the database:
```bash
make init-db
```

The database file is located at `./backend/data/globetrotter.db`.

## Project Structure

```
backend/
├── api/              # API handlers
│   └── handlers.go   # Request handlers
├── cmd/              # Command-line tools
│   └── init_db/      # Database initialization tool
├── data/             # Data files
│   ├── data.json     # Destination data
│   └── globetrotter.db # SQLite database
├── db/               # Database package
│   └── db.go         # Database initialization and operations
├── models/           # Data models
│   └── models.go     # Struct definitions
├── services/         # Business logic
│   ├── destination_service.go # Destination operations
│   └── user_service.go        # User operations
└── main.go           # Entry point
```

## API Endpoints

| Method | Endpoint                    | Description                           |
|--------|----------------------------|---------------------------------------|
| GET    | /api/destinations/random   | Get a random destination with options |
| POST   | /api/users                 | Create a new user                     |
| GET    | /api/users/:username       | Get user information                  |
| POST   | /api/users/:username/score | Update a user's score                 |

## Development

### Running with Hot Reload

For development, you can use Air for hot reloading:

1. Install Air:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. Run with hot reload:
   ```bash
   make dev
   ```

### Running Tests

```bash
make test
```

## Deployment

To build for production:

```bash
make build
```

Then run the compiled binary:

```bash
./bin/globetrotter
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Path to SQLite database file (default: "./backend/data/globetrotter.db")

## License

This project is licensed under the MIT License. 