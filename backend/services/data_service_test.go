package services

import (
	"testing"

	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
	"github.com/stretchr/testify/assert"
)

func setupTestService() *DataService {
	testDB, _ := db.NewDatabase(":memory:")
	return NewDataService(testDB)
}

func TestGetRandomDestination(t *testing.T) {
	service := setupTestService()

	destination, err := service.GetRandomDestination()

	assert.Nil(t, err)
	assert.NotEmpty(t, destination.Name)
	assert.NotEmpty(t, destination.Country)
}

func TestCreateAndGetUser(t *testing.T) {
	service := setupTestService()

	// Create a test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	createdUser, err := service.CreateUser(user)

	assert.Nil(t, err)
	assert.Equal(t, "testuser", createdUser.Username)
	assert.Equal(t, "test@example.com", createdUser.Email)

	// Get the user
	fetchedUser, err := service.GetUser("testuser")

	assert.Nil(t, err)
	assert.Equal(t, "testuser", fetchedUser.Username)
	assert.Equal(t, "test@example.com", fetchedUser.Email)
}

func TestGameOperations(t *testing.T) {
	service := setupTestService()

	// Create a game
	game, err := service.CreateGame("testplayer", "easy")

	assert.Nil(t, err)
	assert.NotEmpty(t, game.ID)
	assert.Equal(t, "testplayer", game.Username)
	assert.Equal(t, "easy", game.Difficulty)

	// Get next question
	question, err := service.GetNextQuestion(game.ID)

	assert.Nil(t, err)
	assert.NotEmpty(t, question.ID)
	assert.NotEmpty(t, question.Text)

	// Submit an answer
	result, err := service.SubmitAnswer(game.ID, "Paris")

	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestGetGameResults(t *testing.T) {
	service := setupTestService()

	// Create a game
	game, err := service.CreateGame("testplayer", "easy")
	assert.Nil(t, err)

	// Get game results
	results, err := service.GetGameResults(game.ID)

	assert.Nil(t, err)
	assert.NotNil(t, results)
}

func TestGetGameSummary(t *testing.T) {
	service := setupTestService()

	// Create a game
	game, err := service.CreateGame("testplayer", "easy")
	assert.Nil(t, err)

	// Get game summary
	summary, err := service.GetGameSummary(game.ID)

	assert.Nil(t, err)
	assert.NotNil(t, summary)
}
