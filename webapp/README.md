# Globetrotter

## The Ultimate Travel Guessing Game

Globetrotter is a full-stack web application where users get cryptic clues about famous places around the world and must guess which destination they refer to. Upon guessing correctly, users unlock fun facts, trivia, and surprises about the destination!

## Features

- 🌍 **Destination Guessing**: Receive cryptic clues about famous places and guess the correct destination
- 🎮 **Interactive Gameplay**: Immediate feedback with animations for correct and incorrect answers
- 🏆 **Score Tracking**: Keep track of your correct and incorrect answers
- 🎲 **Random Destinations**: Play with a diverse set of 100+ destinations from around the world
- 🔗 **Challenge Friends**: Generate unique links to challenge friends to beat your score
- 📱 **Responsive Design**: Play on any device with a fully responsive UI

## Tech Stack

### Frontend
- **React**: UI library for building the user interface
- **Material UI**: Component library for modern, responsive design
- **React Router**: For client-side routing
- **Axios**: For API requests to the backend
- **React Confetti**: For celebration animations

### Backend
- **Go (Golang)**: For the backend API server
- **Gin**: Web framework for routing and middleware
- **SQLite**: Lightweight database for storing user data and destinations
- **Docker**: For containerization and deployment

## Getting Started

### Prerequisites
- Node.js (v14 or higher)
- npm or yarn
- Go (v1.18 or higher) - only needed for backend development

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/globetrotter.git
   cd globetrotter
   ```

2. Install frontend dependencies:
   ```
   cd webapp
   npm install
   ```

3. Start the frontend development server:
   ```
   npm start
   ```

4. In a separate terminal, start the backend server:
   ```
   cd ../backend
   go run main.go
   ```

5. Open your browser and navigate to `http://localhost:3000`

## Docker Deployment

To run the entire application using Docker:

```
docker build -t globetrotter .
docker run -p 8080:8080 globetrotter
```

Then visit `http://localhost:8080` in your browser.

## Project Structure

```
globetrotter/
├── webapp/                # Frontend React application
│   ├── public/            # Static assets
│   ├── src/               # React source code
│   └── package.json       # Frontend dependencies
├── backend/               # Go backend API
│   ├── api/               # HTTP handlers and routes
│   ├── db/                # Database access layer
│   ├── models/            # Data models
│   ├── services/          # Business logic
│   └── main.go            # Application entry point
├── Dockerfile             # Docker configuration
└── README.md              # This file
```
