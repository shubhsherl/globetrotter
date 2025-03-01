package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
	"github.com/shubhsherl/globetrotter/backend/services/images"
)

// GameService handles game-related operations
type GameService struct {
	db *db.Database
}

// NewGameService creates a new game service
func NewGameService(database *db.Database) *GameService {
	return &GameService{
		db: database,
	}
}

// CreateGame creates a new game for a user
func (s *GameService) CreateGame(userID int) (int, error) {
	// Create a new game
	gameID, err := s.db.CreateGame(userID, 5) // 5 questions per game
	if err != nil {
		return 0, err
	}

	// Get all destinations
	destinations, err := s.db.GetAllDestinations()
	if err != nil {
		return 0, err
	}

	// Shuffle destinations
	rand.Shuffle(len(destinations), func(i, j int) {
		destinations[i], destinations[j] = destinations[j], destinations[i]
	})

	// Select 5 destinations for questions
	questionDestinations := destinations[:5]

	// Create 5 questions
	for _, dest := range questionDestinations {
		// Use a random clue as the question
		var question string
		if len(dest.Clues) > 0 {
			// Pick a random clue
			question = dest.Clues[rand.Intn(len(dest.Clues))]
		} else {
			// Fallback to default question if no clues available
			question = fmt.Sprintf("Where is %s located?", dest.City)
		}

		// Generate options (3 wrong options + 1 correct)
		optionDestinations := []models.Destination{dest} // Add correct destination

		// Create a pool of wrong options (excluding the correct answer)
		wrongOptionPool := make([]models.Destination, 0)
		for _, wrongDest := range destinations {
			if wrongDest.ID != dest.ID {
				wrongOptionPool = append(wrongOptionPool, wrongDest)
			}
		}

		// Shuffle the wrong option pool
		rand.Shuffle(len(wrongOptionPool), func(i, j int) {
			wrongOptionPool[i], wrongOptionPool[j] = wrongOptionPool[j], wrongOptionPool[i]
		})

		// Add 3 random wrong options
		optionDestinations = append(optionDestinations, wrongOptionPool[:3]...)

		// Shuffle options
		rand.Shuffle(len(optionDestinations), func(i, j int) {
			optionDestinations[i], optionDestinations[j] = optionDestinations[j], optionDestinations[i]
		})

		// Create option destination IDs
		optionDestinationIDs := make([]int, len(optionDestinations))
		for k, optDest := range optionDestinations {
			optionDestinationIDs[k] = optDest.ID
		}

		// Add question to game
		_, err = s.db.AddGameQuestion(gameID, question, optionDestinationIDs, dest.ID)
		if err != nil {
			return 0, err
		}
	}

	return gameID, nil
}

// GetNextQuestion gets the next unanswered question for a game
func (s *GameService) GetNextQuestion(gameID int) (*models.GameQuestionDetail, error) {
	// Get the next question
	question, err := s.db.GetNextQuestion(gameID)
	if err != nil {
		return nil, err
	}

	// Don't return the correct destination ID to the client
	question.CorrectDestinationID = 0

	return question, nil
}

// SubmitAnswer submits an answer for a question
func (s *GameService) SubmitAnswer(gameID, questionID, selectedDestinationID int) (*models.SubmitAnswerResponse, error) {
	// Check if the question has already been answered
	question, err := s.db.GetQuestionByID(gameID, questionID)
	if err != nil {
		return nil, err
	}

	if question.IsAnswered == 1 {
		return nil, fmt.Errorf("question already answered")
	}

	// Validate that the selected destination ID is in the list of options
	isValidOption := false
	for _, optionID := range question.OptionDestinationIDs {
		if optionID == selectedDestinationID {
			isValidOption = true
			break
		}
	}

	if !isValidOption {
		return nil, fmt.Errorf("selected destination is not in the list of options")
	}

	// Get the correct destination details
	correctDest, err := s.db.GetDestinationByID(question.CorrectDestinationID)
	if err != nil {
		return nil, err
	}

	// Submit the answer
	err = s.db.SubmitAnswer(gameID, questionID, selectedDestinationID)
	if err != nil {
		return nil, err
	}

	// Check if the answer is correct
	isCorrect := selectedDestinationID == question.CorrectDestinationID

	// Prepare response
	response := &models.SubmitAnswerResponse{
		Correct:         isCorrect,
		CorrectCity:     correctDest.City,
		CorrectCountry:  correctDest.Country,
		CorrectOptionID: question.CorrectDestinationID,
	}

	// Add fun fact or trivia based on correctness
	if isCorrect && len(correctDest.FunFact) > 0 {
		// Pick a random fun fact
		response.FunFact = correctDest.FunFact[rand.Intn(len(correctDest.FunFact))]
	} else if !isCorrect && len(correctDest.Trivia) > 0 {
		// Pick a random trivia
		response.Trivia = correctDest.Trivia[rand.Intn(len(correctDest.Trivia))]
	}

	return response, nil
}

// GetGameResult gets the result of a game
func (s *GameService) GetGameResult(gameID int) (*models.GameResult, error) {
	return s.db.GetGameResult(gameID)
}

// HasNextQuestion checks if a game has more unanswered questions
func (s *GameService) HasNextQuestion(gameID int) (bool, error) {
	return s.db.HasNextQuestion(gameID)
}

// GetDestinationByID gets a destination by its ID
func (s *GameService) GetDestinationByID(destinationID int) (*models.Destination, error) {
	return s.db.GetDestinationByID(destinationID)
}

// GetGameSummary gets a summary of a game
func (s *GameService) GetGameSummary(gameID int) (*models.GameSummary, error) {
	// Get the game
	game, err := s.db.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	// Get the user
	user, err := s.db.GetUserByID(game.UserID)
	if err != nil {
		return nil, err
	}

	// Create summary
	summary := &models.GameSummary{
		GameID:         game.ID,
		Username:       user.Username,
		ImageURL:       images.GetTravelImage(),
		TotalQuestions: game.TotalQuestions,
		TotalAnswered:  game.TotalAnswered,
		TotalCorrect:   game.TotalCorrect,
		CreatedAt:      game.CreatedAt.Format(time.RFC3339),
	}

	return summary, nil
}
