package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shubhsherl/globetrotter/backend/models"
)

var DB *sql.DB
var DBx *sqlx.DB

// InitDB initializes the SQLite database
func InitDB(dbPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	// Open database connection
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Initialize sqlx
	DBx = sqlx.NewDb(DB, "sqlite3")

	// Create tables if they don't exist
	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Check if destinations table is empty, if so, seed data
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM destinations").Scan(&count)
	if err != nil || count == 0 {
		if err := seedDestinations(); err != nil {
			return fmt.Errorf("failed to seed destinations: %v", err)
		}
	}

	log.Println("Database initialized successfully")
	return nil
}

// createTables creates the necessary tables in the database
func createTables() error {
	// Create destinations table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS destinations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			city TEXT NOT NULL,
			country TEXT NOT NULL,
			clues TEXT NOT NULL,
			fun_facts TEXT NOT NULL,
			trivia TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create users table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create games table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			total_questions INTEGER DEFAULT 5,
			total_correct INTEGER DEFAULT 0,
			total_incorrect INTEGER DEFAULT 0,
			total_answered INTEGER DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users (id)
		)
	`)
	if err != nil {
		return err
	}

	// Create game_questions table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS game_questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			question TEXT NOT NULL,
			options TEXT NOT NULL,
			correct_destination_id INTEGER NOT NULL,
			selected_destination_id INTEGER DEFAULT 0,
			is_answered INTEGER DEFAULT 0,
			FOREIGN KEY (game_id) REFERENCES games (id),
			FOREIGN KEY (correct_destination_id) REFERENCES destinations (id)
		)
	`)

	return err
}

// seedDestinations loads destination data from JSON file and inserts into database
func seedDestinations() error {
	// Read JSON file
	data, err := ioutil.ReadFile("data/data.json")
	if err != nil {
		return fmt.Errorf("failed to read data.json: %v", err)
	}

	// Parse JSON
	var destinations []models.Destination
	if err := json.Unmarshal(data, &destinations); err != nil {
		return fmt.Errorf("failed to parse data.json: %v", err)
	}

	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.Prepare(`
		INSERT INTO destinations (city, country, clues, fun_facts, trivia)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Insert destinations
	for _, dest := range destinations {
		// Convert slices to JSON strings
		cluesJSON, err := json.Marshal(dest.Clues)
		if err != nil {
			return err
		}

		funFactsJSON, err := json.Marshal(dest.FunFact)
		if err != nil {
			return err
		}

		triviaJSON, err := json.Marshal(dest.Trivia)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(dest.City, dest.Country, cluesJSON, funFactsJSON, triviaJSON)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Seeded %d destinations", len(destinations))
	return nil
}

// Database handles database operations
type Database struct {
	db  *sql.DB
	dbx *sqlx.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, "sqlite3")

	return &Database{db: db, dbx: dbx}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetAllDestinations retrieves all destinations from the database
func (d *Database) GetAllDestinations() ([]models.Destination, error) {
	type DestinationWithJSON struct {
		ID           int    `db:"id"`
		City         string `db:"city"`
		Country      string `db:"country"`
		CluesJSON    string `db:"clues"`
		FunFactsJSON string `db:"fun_facts"`
		TriviaJSON   string `db:"trivia"`
	}

	var destinationsWithJSON []DestinationWithJSON

	err := d.dbx.Select(&destinationsWithJSON, `
		SELECT id, city, country, clues, fun_facts, trivia
		FROM destinations
	`)

	if err != nil {
		return nil, err
	}

	var destinations []models.Destination
	for _, d := range destinationsWithJSON {
		dest := models.Destination{
			ID:      d.ID,
			City:    d.City,
			Country: d.Country,
		}

		// Parse JSON strings
		if err := json.Unmarshal([]byte(d.CluesJSON), &dest.Clues); err != nil {
			return nil, fmt.Errorf("failed to parse clues: %v", err)
		}

		if err := json.Unmarshal([]byte(d.FunFactsJSON), &dest.FunFact); err != nil {
			return nil, fmt.Errorf("failed to parse fun facts: %v", err)
		}

		if err := json.Unmarshal([]byte(d.TriviaJSON), &dest.Trivia); err != nil {
			return nil, fmt.Errorf("failed to parse trivia: %v", err)
		}

		destinations = append(destinations, dest)
	}

	return destinations, nil
}

// GetUserByUsername retrieves a user by username
func (d *Database) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := d.dbx.Get(&user, `
		SELECT id, username, created_at
		FROM users
		WHERE username = ?
	`, username)

	return user, err
}

// SaveUser saves a user to the database
func (d *Database) SaveUser(user models.User) error {
	// Check if user exists
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Update existing user
		_, err = d.db.Exec("UPDATE users SET created_at = ? WHERE username = ?",
			user.CreatedAt, user.Username)
	} else {
		// Insert new user
		_, err = d.db.Exec("INSERT INTO users (username, created_at) VALUES (?, ?)",
			user.Username, user.CreatedAt)
	}

	return err
}

// GetDB returns the database instance
func GetDB() *Database {
	return &Database{db: DB, dbx: DBx}
}

// CreateGame creates a new game for a user
func (d *Database) CreateGame(userID int, totalQuestions int) (int, error) {
	result, err := d.db.Exec(`
		INSERT INTO games (user_id, total_questions)
		VALUES (?, ?)
	`, userID, totalQuestions)
	if err != nil {
		return 0, err
	}

	gameID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(gameID), nil
}

// AddGameQuestion adds a question to a game
func (d *Database) AddGameQuestion(gameID int, question string, optionDestinationIDs []int, correctDestinationID int) (int, error) {
	// Convert options to JSON
	optionsJSON, err := json.Marshal(optionDestinationIDs)
	if err != nil {
		return 0, err
	}

	result, err := d.db.Exec(`
		INSERT INTO game_questions (game_id, question, options, correct_destination_id)
		VALUES (?, ?, ?, ?)
	`, gameID, question, string(optionsJSON), correctDestinationID)
	if err != nil {
		return 0, err
	}

	questionID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(questionID), nil
}

// GetNextQuestion gets the next unanswered question for a game
func (d *Database) GetNextQuestion(gameID int) (*models.GameQuestionDetail, error) {
	type QuestionWithOptionsJSON struct {
		models.GameQuestionDetail
		OptionsJSON string `db:"options_json"`
	}

	var questionWithJSON QuestionWithOptionsJSON

	err := d.dbx.Get(&questionWithJSON, `
		SELECT id, game_id, question, options as options_json, 
		       correct_destination_id, selected_destination_id, is_answered
		FROM game_questions
		WHERE game_id = ? AND is_answered = 0
		ORDER BY id ASC
		LIMIT 1
	`, gameID)

	if err != nil {
		return nil, err
	}

	// Parse options JSON
	if err := json.Unmarshal([]byte(questionWithJSON.OptionsJSON), &questionWithJSON.OptionDestinationIDs); err != nil {
		return nil, err
	}

	return &questionWithJSON.GameQuestionDetail, nil
}

// GetQuestionByID gets a question by its ID
func (d *Database) GetQuestionByID(gameID, questionID int) (*models.GameQuestionDetail, error) {
	type QuestionWithOptionsJSON struct {
		models.GameQuestionDetail
		OptionsJSON string `db:"options_json"`
	}

	var questionWithJSON QuestionWithOptionsJSON

	err := d.dbx.Get(&questionWithJSON, `
		SELECT id, game_id, question, options as options_json, 
		       correct_destination_id, selected_destination_id, 
		       is_answered
		FROM game_questions
		WHERE game_id = ? AND id = ?
	`, gameID, questionID)

	if err != nil {
		return nil, err
	}

	// Parse options JSON
	if err := json.Unmarshal([]byte(questionWithJSON.OptionsJSON), &questionWithJSON.OptionDestinationIDs); err != nil {
		return nil, err
	}

	return &questionWithJSON.GameQuestionDetail, nil
}

// SubmitAnswer submits an answer for a question
func (d *Database) SubmitAnswer(gameID, questionID int, selectedDestinationID int) error {
	// Update the question
	_, err := d.db.Exec(`
		UPDATE game_questions
		SET selected_destination_id = ?, is_answered = 1
		WHERE id = ? AND game_id = ?
	`, selectedDestinationID, questionID, gameID)

	if err != nil {
		return err
	}

	// Get the correct destination ID
	var correctDestinationID int
	err = d.db.QueryRow(`
		SELECT correct_destination_id
		FROM game_questions
		WHERE id = ? AND game_id = ?
	`, questionID, gameID).Scan(&correctDestinationID)

	if err != nil {
		return err
	}

	// Check if the answer is correct
	isCorrect := selectedDestinationID == correctDestinationID

	// Update the game stats
	if isCorrect {
		_, err = d.db.Exec(`
			UPDATE games
			SET total_correct = total_correct + 1, total_answered = total_answered + 1
			WHERE id = ?
		`, gameID)
	} else {
		_, err = d.db.Exec(`
			UPDATE games
			SET total_incorrect = total_incorrect + 1, total_answered = total_answered + 1
			WHERE id = ?
		`, gameID)
	}

	return err
}

// GetGameResult gets the result of a game
func (d *Database) GetGameResult(gameID int) (*models.GameResult, error) {
	// Get game info
	var game models.Game
	err := d.dbx.Get(&game, `
		SELECT id, user_id, total_questions, 
		       total_correct, total_incorrect, 
		       total_answered, created_at
		FROM games
		WHERE id = ?
	`, gameID)

	if err != nil {
		return nil, err
	}

	// Get questions
	type QuestionWithOptionsJSON struct {
		models.GameQuestionDetail
		OptionsJSON string `db:"options_json"`
	}

	var questionsWithJSON []QuestionWithOptionsJSON

	err = d.dbx.Select(&questionsWithJSON, `
		SELECT id, game_id, question, options as options_json, 
		       correct_destination_id, selected_destination_id, 
		       is_answered
		FROM game_questions
		WHERE game_id = ?
	`, gameID)

	if err != nil {
		return nil, err
	}

	// Parse options JSON for each question
	var questions []models.GameQuestionDetail
	for _, q := range questionsWithJSON {
		if err := json.Unmarshal([]byte(q.OptionsJSON), &q.OptionDestinationIDs); err != nil {
			return nil, err
		}

		if q.IsAnswered == 0 {
			q.CorrectDestinationID = 0
		}

		questions = append(questions, q.GameQuestionDetail)
	}

	return &models.GameResult{
		GameID:         game.ID,
		TotalQuestions: game.TotalQuestions,
		TotalCorrect:   game.TotalCorrect,
		TotalIncorrect: game.TotalIncorrect,
		Questions:      questions,
	}, nil
}

// HasNextQuestion checks if a game has more unanswered questions
func (d *Database) HasNextQuestion(gameID int) (bool, error) {
	var count int
	err := d.db.QueryRow(`
		SELECT COUNT(*)
		FROM game_questions
		WHERE game_id = ? AND is_answered = 0
	`, gameID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 1, nil
}

// GetUserIDByUsername gets a user's ID by username
func (d *Database) GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := d.db.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?
	`, username).Scan(&userID)

	return userID, err
}

// GetDestinationByID gets a destination by its ID
func (d *Database) GetDestinationByID(destinationID int) (*models.Destination, error) {
	type DestinationWithJSON struct {
		ID           int    `db:"id"`
		City         string `db:"city"`
		Country      string `db:"country"`
		CluesJSON    string `db:"clues"`
		FunFactsJSON string `db:"fun_facts"`
		TriviaJSON   string `db:"trivia"`
	}

	var destWithJSON DestinationWithJSON

	err := d.dbx.Get(&destWithJSON, `
		SELECT id, city, country, clues, fun_facts, trivia
		FROM destinations
		WHERE id = ?
	`, destinationID)

	if err != nil {
		return nil, err
	}

	dest := &models.Destination{
		ID:      destWithJSON.ID,
		City:    destWithJSON.City,
		Country: destWithJSON.Country,
	}

	// Parse JSON strings
	if err := json.Unmarshal([]byte(destWithJSON.CluesJSON), &dest.Clues); err != nil {
		return nil, fmt.Errorf("failed to parse clues: %v", err)
	}

	if err := json.Unmarshal([]byte(destWithJSON.FunFactsJSON), &dest.FunFact); err != nil {
		return nil, fmt.Errorf("failed to parse fun facts: %v", err)
	}

	if err := json.Unmarshal([]byte(destWithJSON.TriviaJSON), &dest.Trivia); err != nil {
		return nil, fmt.Errorf("failed to parse trivia: %v", err)
	}

	return dest, nil
}

// GetGame gets a game by ID
func (d *Database) GetGame(gameID int) (*models.Game, error) {
	var game models.Game
	err := d.dbx.Get(&game, `
		SELECT id, user_id, total_questions, 
		       total_correct, total_incorrect, 
		       total_answered, created_at
		FROM games
		WHERE id = ?
	`, gameID)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

// GetUserByID gets a user by ID
func (d *Database) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	err := d.dbx.Get(&user, `
		SELECT id, username, created_at
		FROM users
		WHERE id = ?
	`, userID)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
