package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Initialize test database
	testDB, _ := db.NewDatabase(":memory:")
	InitServices(testDB)

	// Setup routes
	SetupRoutes(r)

	return r
}

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Contains(t, response, "uptime")
}

func TestGetRandomDestination(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/destinations/random", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Contains(t, response, "destination")
}

func TestCreateUser(t *testing.T) {
	router := setupTestRouter()

	// Create a test user
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
	}

	jsonValue, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, "testuser", response["username"])
}

func TestGetUser(t *testing.T) {
	router := setupTestRouter()

	// First create a user
	user := models.User{
		Username: "testuser2",
		Email:    "test2@example.com",
	}

	jsonValue, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Now test getting the user
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/users/testuser2", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, "testuser2", response["username"])
}

func TestGameFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. Start a game
	gameRequest := map[string]interface{}{
		"username":   "testplayer",
		"difficulty": "easy",
	}

	jsonValue, _ := json.Marshal(gameRequest)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/game/play", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var gameResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &gameResponse)

	assert.Nil(t, err)
	assert.Contains(t, gameResponse, "game_id")

	gameID := gameResponse["game_id"].(string)

	// 2. Get next question
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/game/"+gameID+"/next-question", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var questionResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &questionResponse)

	assert.Nil(t, err)
	assert.Contains(t, questionResponse, "question")

	// 3. Submit an answer
	answerRequest := map[string]interface{}{
		"answer": "Paris",
	}

	jsonValue, _ = json.Marshal(answerRequest)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/game/"+gameID+"/submit-answer", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var answerResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &answerResponse)

	assert.Nil(t, err)
	assert.Contains(t, answerResponse, "correct")
}
