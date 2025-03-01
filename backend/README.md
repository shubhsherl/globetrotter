# Globetrotter Challenge Backend

This is the backend API for the Globetrotter Challenge game, built with Go (Golang), the Gin web framework, and SQLite database.

## Features

- RESTful API for game data and user management
- Random destination selection with multiple-choice options
- Challenge sharing functionality with social media meta tags
- SQLite database for persistent storage
- Integration with Pexels API for destination images

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

3. Configure environment variables:
   ```bash
   cp .env.example .env
   ```
   
   Then edit the `.env` file to add your Pexels API key.

4. Build the application:
   ```bash
   make build
   ```

5. Run the server:
   ```bash
   make run
   ```

   The server will start on port 8080 by default.

## Database Management

The application uses SQLite for data storage. The database will be automatically initialized when you run the server for the first time.

To run database migrations:
```bash
make migrate
```

To reinitialize the database:
```bash
make init-db
```

The database file is located at `./data/globetrotter.db`.

## Project Structure

```
backend/
├── api/              # API handlers
│   └── handlers.go   # Request handlers
├── cmd/              # Command-line tools
│   ├── init_db/      # Database initialization tool
│   └── migrate/      # Database migration tool
├── data/             # Data files
│   └── globetrotter.db # SQLite database
├── db/               # Database package
│   └── db.go         # Database initialization and operations
├── migrations/       # SQL migration files
│   ├── 001_initial_schema.sql
│   ├── 002_add_migrations_table.sql
│   └── 003_add_indexes.sql
├── models/           # Data models
│   └── models.go     # Struct definitions
├── services/         # Business logic
│   ├── data_service.go      # Data operations
│   ├── destination_service.go # Destination operations
│   ├── game_service.go      # Game operations
│   ├── user_service.go      # User operations
│   └── images/             # Image service
├── .env              # Environment variables
├── .env.example      # Example environment variables
├── Makefile          # Build and run commands
└── main.go           # Entry point
```

## API Endpoints

| Method | Endpoint                    | Description                           |
|--------|----------------------------|---------------------------------------|
| GET    | /health                    | Health check endpoint                 |
| GET    | /api/destinations/random   | Get a random destination              |
| POST   | /api/users                 | Create a new user                     |
| GET    | /api/users/:username       | Get user information                  |
| POST   | /api/game/play             | Start a new game                      |
| GET    | /api/game/:id/next-question| Get the next question in a game       |
| POST   | /api/game/:id/submit-answer| Submit an answer for a question       |
| GET    | /api/game/:id/result       | Get the result of a game              |
| GET    | /api/game/:id/summary      | Get a summary of a completed game     |
| GET    | /challenge/:username       | Serve challenge page with meta tags   |
| GET    | /challenge/:username/:gameID | Serve specific game challenge page  |

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
- `DB_PATH`: Path to SQLite database file (default: "./data/globetrotter.db")
- `PEXELS_API_KEY`: API key for Pexels image service

## License

This project is licensed under the MIT License. 