package db

import (
	"testing"

	"github.com/shubhsherl/globetrotter/backend/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *Database {
	db, _ := NewDatabase(":memory:")
	return db
}

func TestDatabaseConnection(t *testing.T) {
	db := setupTestDB()

	// Test that the database connection is valid
	err := db.DB.Ping()
	assert.Nil(t, err)
}

func TestUserOperations(t *testing.T) {
	db := setupTestDB()

	// Create a test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Insert the user
	err := db.CreateUser(&user)
	assert.Nil(t, err)

	// Get the user
	fetchedUser, err := db.GetUser("testuser")
	assert.Nil(t, err)
	assert.Equal(t, "testuser", fetchedUser.Username)
	assert.Equal(t, "test@example.com", fetchedUser.Email)
}

func TestDestinationOperations(t *testing.T) {
	db := setupTestDB()

	// Get a random destination
	destination, err := db.GetRandomDestination()
	assert.Nil(t, err)
	assert.NotEmpty(t, destination.Name)
}

func TestGameOperations(t *testing.T) {
	db := setupTestDB()

	// Create a game
	game := models.Game{
		ID:         "game123",
		Username:   "testplayer",
		Difficulty: "easy",
		Status:     "in_progress",
	}

	// Insert the game
	err := db.CreateGame(&game)
	assert.Nil(t, err)

	// Get the game
	fetchedGame, err := db.GetGame("game123")
	assert.Nil(t, err)
	assert.Equal(t, "game123", fetchedGame.ID)
	assert.Equal(t, "testplayer", fetchedGame.Username)
	assert.Equal(t, "easy", fetchedGame.Difficulty)
}

func TestQuestionOperations(t *testing.T) {
	db := setupTestDB()

	// Create a game first
	game := models.Game{
		ID:         "game456",
		Username:   "testplayer",
		Difficulty: "easy",
		Status:     "in_progress",
	}

	err := db.CreateGame(&game)
	assert.Nil(t, err)

	// Create a question
	question := models.Question{
		GameID:        "game456",
		Text:          "What is the capital of France?",
		CorrectAnswer: "Paris",
		Options:       []string{"London", "Paris", "Berlin", "Madrid"},
		Type:          "multiple_choice",
	}

	// Insert the question
	err = db.CreateQuestion(&question)
	assert.Nil(t, err)
	assert.NotZero(t, question.ID)

	// Get the question
	fetchedQuestion, err := db.GetQuestion(question.ID)
	assert.Nil(t, err)
	assert.Equal(t, "What is the capital of France?", fetchedQuestion.Text)
	assert.Equal(t, "Paris", fetchedQuestion.CorrectAnswer)
}

func TestAnswerOperations(t *testing.T) {
	db := setupTestDB()

	// Create a game first
	game := models.Game{
		ID:         "game789",
		Username:   "testplayer",
		Difficulty: "easy",
		Status:     "in_progress",
	}

	err := db.CreateGame(&game)
	assert.Nil(t, err)

	// Create a question
	question := models.Question{
		GameID:        "game789",
		Text:          "What is the capital of France?",
		CorrectAnswer: "Paris",
		Options:       []string{"London", "Paris", "Berlin", "Madrid"},
		Type:          "multiple_choice",
	}

	err = db.CreateQuestion(&question)
	assert.Nil(t, err)

	// Create an answer
	answer := models.Answer{
		QuestionID: question.ID,
		GameID:     "game789",
		UserAnswer: "Paris",
		IsCorrect:  true,
		TimeTaken:  5.2,
	}

	// Insert the answer
	err = db.CreateAnswer(&answer)
	assert.Nil(t, err)
	assert.NotZero(t, answer.ID)

	// Get the answer
	fetchedAnswer, err := db.GetAnswer(answer.ID)
	assert.Nil(t, err)
	assert.Equal(t, "Paris", fetchedAnswer.UserAnswer)
	assert.True(t, fetchedAnswer.IsCorrect)
}
