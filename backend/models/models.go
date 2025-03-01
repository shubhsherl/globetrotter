package models

import (
	"encoding/json"
	"math/rand"
	"time"
)

// Destination represents a location in the game
type Destination struct {
	ID      int      `json:"id,omitempty" db:"id"`
	City    string   `json:"city" db:"city"`
	Country string   `json:"country" db:"country"`
	Clues   []string `json:"clues" db:"-"`
	FunFact []string `json:"fun_fact" db:"-"`
	Trivia  []string `json:"trivia" db:"-"`
}

// User represents a player in the game
type User struct {
	ID        int    `json:"id,omitempty" db:"id"`
	Username  string `json:"username" db:"username"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
}

// GameQuestion represents a question for the frontend
type GameQuestion struct {
	ID            int      `json:"id,omitempty" db:"id"`
	Clues         []string `json:"clues" db:"-"`
	Options       []string `json:"options" db:"-"`
	CorrectIndex  int      `json:"-" db:"correct_index"` // Not sent to client
	CorrectAnswer string   `json:"correct_answer" db:"correct_answer"`
}

// Game represents a game session
type Game struct {
	ID             int       `json:"id,omitempty" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	TotalQuestions int       `json:"total_questions" db:"total_questions"`
	TotalCorrect   int       `json:"total_correct" db:"total_correct"`
	TotalIncorrect int       `json:"total_incorrect" db:"total_incorrect"`
	TotalAnswered  int       `json:"total_answered" db:"total_answered"`
}

// GameQuestionDetail represents a question in a game
type GameQuestionDetail struct {
	ID                    int    `json:"id,omitempty" db:"id"`
	GameID                int    `json:"game_id" db:"game_id"`
	Question              string `json:"question" db:"question"`
	OptionDestinationIDs  []int  `json:"options" db:"-"`                                               // Changed from []string to []int to store destination IDs
	CorrectDestinationID  int    `json:"correct_destination_id,omitempty" db:"correct_destination_id"` // Not sent to client during game
	SelectedDestinationID int    `json:"selected_destination_id,omitempty" db:"selected_destination_id"`
	IsAnswered            int    `json:"is_answered" db:"is_answered"` // 0 = false, 1 = true
}

// NextQuestionResponse represents the response for the next question API
type NextQuestionResponse struct {
	GameID         int            `json:"game_id" db:"game_id"`
	QuestionID     int            `json:"question_id" db:"question_id"`
	Question       string         `json:"question" db:"question"`
	Options        []string       `json:"options" db:"-"`
	OptionsDisplay map[int]string `json:"options_display" db:"-"`
	HasNext        bool           `json:"has_next" db:"has_next"`
}

// SubmitAnswerRequest represents the request for submitting an answer
type SubmitAnswerRequest struct {
	GameID              int `json:"game_id" binding:"required" db:"game_id"`
	QuestionID          int `json:"question_id" binding:"required" db:"question_id"`
	SelectedDestination int `json:"selected_destination" binding:"required" db:"selected_destination"`
}

// SubmitAnswerResponse represents the response for submitting an answer
type SubmitAnswerResponse struct {
	Correct        bool   `json:"correct" db:"correct"`
	FunFact        string `json:"fun_fact,omitempty" db:"fun_fact"` // Sent when answer is correct
	Trivia         string `json:"trivia,omitempty" db:"trivia"`     // Sent when answer is incorrect
	CorrectCity    string `json:"correct_city" db:"correct_city"`
	CorrectCountry string `json:"correct_country" db:"correct_country"`
}

// GameResult represents the result of a completed game
type GameResult struct {
	GameID         int                  `json:"game_id" db:"game_id"`
	TotalQuestions int                  `json:"total_questions" db:"total_questions"`
	TotalCorrect   int                  `json:"total_correct" db:"total_correct"`
	TotalIncorrect int                  `json:"total_incorrect" db:"total_incorrect"`
	Questions      []GameQuestionDetail `json:"questions" db:"-"`
}

// DBDestination is used for database operations
type DBDestination struct {
	ID       int    `db:"id"`
	City     string `db:"city"`
	Country  string `db:"country"`
	Clues    string `db:"clues"`     // JSON string
	FunFacts string `db:"fun_facts"` // JSON string
	Trivia   string `db:"trivia"`    // JSON string
}

// ToDestination converts a DBDestination to a Destination
func (d *DBDestination) ToDestination() (*Destination, error) {
	dest := &Destination{
		ID:      d.ID,
		City:    d.City,
		Country: d.Country,
	}

	// Unmarshal JSON strings to slices
	if err := json.Unmarshal([]byte(d.Clues), &dest.Clues); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(d.FunFacts), &dest.FunFact); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(d.Trivia), &dest.Trivia); err != nil {
		return nil, err
	}

	return dest, nil
}

// GenerateOptions creates multiple choice options for a destination
func GenerateOptions(correctDest *Destination, allDests []*Destination) []string {
	// Create correct answer
	correctAnswer := correctDest.City + ", " + correctDest.Country

	// Create a set of options including the correct answer
	options := []string{correctAnswer}

	// Add 3 random incorrect options
	for len(options) < 4 {
		// Pick a random destination
		randIndex := rand.Intn(len(allDests))
		randDest := allDests[randIndex]

		// Skip if it's the correct destination
		if randDest.ID == correctDest.ID {
			continue
		}

		// Create option
		option := randDest.City + ", " + randDest.Country

		// Check if option is already in options
		duplicate := false
		for _, existing := range options {
			if existing == option {
				duplicate = true
				break
			}
		}

		// Add if not duplicate
		if !duplicate {
			options = append(options, option)
		}
	}

	// Shuffle options
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

// GameSummary represents a summary of a game
type GameSummary struct {
	GameID         int    `json:"game_id" db:"game_id"`
	Username       string `json:"username" db:"username"`
	TotalQuestions int    `json:"total_questions" db:"total_questions"`
	TotalAnswered  int    `json:"total_answered" db:"total_answered"`
	TotalCorrect   int    `json:"total_correct" db:"total_correct"`
	CreatedAt      string `json:"created_at" db:"created_at"`
}
