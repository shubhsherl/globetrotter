package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	// Test user creation
	user := User{
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}
	
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotZero(t, user.CreatedAt)
}

func TestDestinationModel(t *testing.T) {
	// Test destination creation
	destination := Destination{
		ID:          1,
		Name:        "Paris",
		Country:     "France",
		Description: "City of Light",
		Latitude:    48.8566,
		Longitude:   2.3522,
		ImageURL:    "paris.jpg",
	}
	
	assert.Equal(t, 1, destination.ID)
	assert.Equal(t, "Paris", destination.Name)
	assert.Equal(t, "France", destination.Country)
	assert.Equal(t, "City of Light", destination.Description)
	assert.Equal(t, 48.8566, destination.Latitude)
	assert.Equal(t, 2.3522, destination.Longitude)
	assert.Equal(t, "paris.jpg", destination.ImageURL)
}

func TestGameModel(t *testing.T) {
	// Test game creation
	game := Game{
		ID:         "game123",
		Username:   "player1",
		Difficulty: "medium",
		StartTime:  time.Now(),
		Status:     "in_progress",
	}
	
	assert.Equal(t, "game123", game.ID)
	assert.Equal(t, "player1", game.Username)
	assert.Equal(t, "medium", game.Difficulty)
	assert.Equal(t, "in_progress", game.Status)
	assert.NotZero(t, game.StartTime)
}

func TestQuestionModel(t *testing.T) {
	// Test question creation
	question := Question{
		ID:            1,
		GameID:        "game123",
		Text:          "What is the capital of France?",
		CorrectAnswer: "Paris",
		Options:       []string{"London", "Paris", "Berlin", "Madrid"},
		Type:          "multiple_choice",
	}
	
	assert.Equal(t, 1, question.ID)
	assert.Equal(t, "game123", question.GameID)
	assert.Equal(t, "What is the capital of France?", question.Text)
	assert.Equal(t, "Paris", question.CorrectAnswer)
	assert.Len(t, question.Options, 4)
	assert.Equal(t, "multiple_choice", question.Type)
}

func TestAnswerModel(t *testing.T) {
	// Test answer creation
	answer := Answer{
		ID:          1,
		QuestionID:  1,
		GameID:      "game123",
		UserAnswer:  "Paris",
		IsCorrect:   true,
		AnsweredAt:  time.Now(),
		TimeTaken:   10.5,
	}
	
	assert.Equal(t, 1, answer.ID)
	assert.Equal(t, 1, answer.QuestionID)
	assert.Equal(t, "game123", answer.GameID)
	assert.Equal(t, "Paris", answer.UserAnswer)
	assert.True(t, answer.IsCorrect)
	assert.NotZero(t, answer.AnsweredAt)
	assert.Equal(t, 10.5, answer.TimeTaken)
} 