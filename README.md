# Globetrotter Challenge

The Globetrotter Challenge is an interactive travel guessing game where players are presented with cryptic clues about famous destinations and must guess the correct location. Upon guessing correctly, players unlock fun facts and trivia about the destination!

## Live Demo

The application is deployed and available at:
**[https://globetrotter.up.railway.app](https://globetrotter.up.railway.app)**

Try it out and test your geography knowledge!

## Project Structure

This project is organized as a monorepo with two main components:

- `backend/`: Go (Golang) backend API
- `webapp/`: React frontend application

## Features

- 🌍 100+ destinations with unique clues, fun facts, and trivia
- 🎮 Interactive gameplay with immediate feedback
- 🎯 Score tracking to monitor your progress
- 🔗 Challenge friends via shareable links
- 🎉 Celebratory animations for correct answers

## Getting Started

### Prerequisites

- Go 1.16+ (for backend)
- Node.js 14+ and npm (for frontend)
- Git
- (Optional) Pexels API Key

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/globetrotter.git
   cd globetrotter
   ```

2. Set up and run the backend:
   ```bash
   cd backend
   cp .env.example .env # Add your Pexels API Key to the .env file
   go mod download
   go run main.go
   ```

3. Set up and run the frontend:
   ```bash
   cd webapp
   npm install
   npm start
   ```

4. Open your browser and navigate to `http://localhost:3000`

For more detailed instructions, see the README files in the respective directories:
- [Backend README](./backend/README.md)
- [Frontend README](./webapp/README.md)

## How to Play

1. Enter your username to start the game
2. Read the clues about a mystery destination
3. Select your answer from the multiple choices
4. Get immediate feedback and learn fun facts
5. Challenge your friends to beat your score!

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Data sources for destination information
- React and Go communities for excellent documentation
- All contributors to this project 