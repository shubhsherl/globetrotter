# Globetrotter Technical Documentation

## Architecture Overview

The Globetrotter Challenge is built using a modern client-server architecture with a clear separation of concerns:

- **Frontend**: React-based single-page application (SPA)
- **Backend**: Go (Golang) REST API with SQLite database

## Backend Architecture

### Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Database**: SQLite with sqlx for enhanced query capabilities
- **Development Tools**: Air (for hot reloading)

### Directory Structure

```
backend/
├── api/          # HTTP handlers and route definitions
├── cmd/          # Command-line tools for DB initialization and migrations
├── data/         # Data files including destination information
├── db/           # Database access layer
├── migrations/   # SQL migration scripts
├── models/       # Data models and structures
├── services/     # Business logic layer
└── main.go       # Application entry point
```

### Key Components

#### Data Models (`models/`)

The application uses several core data models:

- **User**: Represents a player with username and score tracking
- **Destination**: Contains information about travel destinations including city, country, clues, fun facts, and trivia
- **Game**: Tracks an active game session with status and player information
- **Question**: Represents a question in a game with options and correct answer
- **Answer**: Records a player's answer to a question

#### Database Layer (`db/`)

The database layer handles:

- Database initialization and connection management
- Table creation and schema management
- CRUD operations for all data models
- Transaction management

The SQLite database is used for simplicity and ease of deployment, with tables for users, destinations, games, and game questions/answers.

#### Services Layer (`services/`)

The services layer implements the business logic:

- **DataService**: Coordinates access to other services
- **DestinationService**: Manages destination data
- **UserService**: Handles user creation and retrieval
- **GameService**: Manages game flow, questions, and scoring

This layer acts as an intermediary between the API handlers and the database, ensuring proper data validation and business rule enforcement.

#### API Layer (`api/`)

The API layer provides RESTful endpoints:

- User management (`/api/users`)
- Game flow (`/api/game/*`)
- Destination information (`/api/destinations/*`)

Each endpoint is implemented as a Gin handler function that processes requests, interacts with the services layer, and returns appropriate responses.

### Data Flow

1. Client makes a request to an API endpoint
2. Gin router directs the request to the appropriate handler
3. Handler validates the request and calls the relevant service method
4. Service implements business logic and interacts with the database
5. Database executes queries and returns results
6. Results flow back through the service and handler to the client

### Error Handling

The backend implements consistent error handling:

- Database errors are captured and translated to appropriate HTTP status codes
- Input validation occurs at both the API and service layers
- Detailed error messages are logged but sanitized before being sent to clients

## Frontend Architecture

### Technology Stack

- **Framework**: React
- **State Management**: React Context API
- **Routing**: React Router
- **Styling**: CSS Modules
- **HTTP Client**: Axios

### Directory Structure

```
webapp/
├── public/       # Static assets
├── src/
│   ├── components/  # Reusable UI components
│   ├── context/     # React context providers
│   ├── pages/       # Page components
│   ├── services/    # API service functions
│   ├── styles/      # Global styles
│   ├── utils/       # Utility functions
│   └── App.js       # Main application component
└── package.json     # Dependencies and scripts
```

### Key Components

#### Game Context (`context/GameContext.js`)

The GameContext provides global state management for:

- User information and score
- Current game state
- Question and answer tracking
- Local storage persistence

#### API Services (`services/api.js`)

The API service layer handles communication with the backend:

- User creation and retrieval
- Game flow (start game, get questions, submit answers)
- Score management

#### Game Component (`pages/Game.js`)

The Game component is the core of the application, managing:

- Game state transitions (initial → playing → finished)
- Question display and answer handling
- Score updates and persistence
- Error handling and recovery

### Local Storage Strategy

The application uses localStorage to persist:

- User information (username, score)
- Current game state
- Timestamps for session management

This allows users to refresh the page or return later without losing their progress.

## Key Implementation Details

### Score Management

User scores are managed locally in the frontend:

1. The `resetUserScore` function in `api.js` handles local score resets
2. The `GameContext` updates localStorage when user data changes
3. The `Game.js` component ensures proper score tracking during gameplay

### Game Flow

The game flow follows these states:

1. **Initial**: User enters username or returns to start screen
2. **Playing**: User is presented with questions and submits answers
3. **Finished**: Game is complete, results are displayed

State transitions are managed in the `Game.js` component with appropriate API calls and context updates.

### Error Handling and Recovery

The application implements several strategies for error handling:

- API call failures are caught and logged
- Local fallbacks are used when server operations fail
- The UI provides feedback for error conditions
- Game state is preserved to prevent data loss

## Deployment Considerations

### Backend Deployment

The Go backend can be deployed as:

- A standalone binary on any supported platform
- A containerized application using Docker
- A cloud service (AWS, GCP, Azure)

### Frontend Deployment

The React frontend can be deployed:

- As static files on any web server
- On CDN services like Netlify or Vercel
- In a containerized environment alongside the backend

### Database Considerations

The SQLite database is file-based and requires:

- Proper backup procedures
- Consideration of concurrent access limitations
- Potential migration to a more robust database for high-load scenarios

## Future Enhancements

Potential areas for enhancement include:

1. **Authentication**: Implement proper user authentication
2. **Multiplayer**: Add real-time multiplayer capabilities
3. **Leaderboards**: Implement global and friend-based leaderboards
4. **Content Management**: Add admin interface for destination management
5. **Performance Optimization**: Implement caching and query optimization
6. **Mobile App**: Develop native mobile applications

## Development Workflow

The project uses Make for common development tasks:

- `make setup`: Set up dependencies
- `make build`: Build both frontend and backend
- `make run`: Run the application
- `make dev`: Run in development mode with hot reloading
- `make clean`: Clean build artifacts 