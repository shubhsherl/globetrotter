package services

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
)

var destinations []models.Destination

func init() {
	rand.Seed(time.Now().UnixNano())
	loadDestinations()
}

func loadDestinations() {
	data, err := ioutil.ReadFile("data/data.json")
	if err != nil {
		panic("Failed to load destinations data: " + err.Error())
	}

	if err := json.Unmarshal(data, &destinations); err != nil {
		panic("Failed to parse destinations data: " + err.Error())
	}
}

// DataService coordinates access to other services
type DataService struct {
	destinationService *DestinationService
	userService        *UserService
	gameService        *GameService
}

// NewDataService creates a new data service
func NewDataService(database *db.Database) *DataService {
	return &DataService{
		destinationService: NewDestinationService(database),
		userService:        NewUserService(database),
		gameService:        NewGameService(database),
	}
}

// Then add methods that delegate to the appropriate service
// For example:
func (s *DataService) GetRandomDestination() (models.Destination, error) {
	return s.destinationService.GetRandomDestination()
}

// CreateUser delegates to the user service
func (s *DataService) CreateUser(username string) (models.User, error) {
	return s.userService.CreateUser(username)
}

// GetUser delegates to the user service
func (s *DataService) GetUser(username string) (models.User, error) {
	return s.userService.GetUser(username)
}

// CreateGame delegates to the game service
func (s *DataService) CreateGame(username string) (int, error) {
	// Get user ID from username
	user, err := s.GetUser(username)
	if err != nil {
		return 0, err
	}
	return s.gameService.CreateGame(user.ID)
}

// GetNextQuestion delegates to the game service
func (s *DataService) GetNextQuestion(gameID int) (*models.GameQuestionDetail, error) {
	return s.gameService.GetNextQuestion(gameID)
}

// SubmitAnswer delegates to the game service
func (s *DataService) SubmitAnswer(gameID, questionID int, selectedDestinationID int) (*models.SubmitAnswerResponse, error) {
	return s.gameService.SubmitAnswer(gameID, questionID, selectedDestinationID)
}

// GetGameResult delegates to the game service
func (s *DataService) GetGameResult(gameID int) (*models.GameResult, error) {
	return s.gameService.GetGameResult(gameID)
}

// HasNextQuestion checks if a game has more unanswered questions
func (s *DataService) HasNextQuestion(gameID int) (bool, error) {
	return s.gameService.HasNextQuestion(gameID)
}

// GetDestinationByID gets a destination by its ID
func (s *DataService) GetDestinationByID(destinationID int) (*models.Destination, error) {
	return s.gameService.GetDestinationByID(destinationID)
}

// GetGameSummary delegates to the game service
func (s *DataService) GetGameSummary(gameID int) (*models.GameSummary, error) {
	return s.gameService.GetGameSummary(gameID)
}
