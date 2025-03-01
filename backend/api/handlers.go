package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
	"github.com/shubhsherl/globetrotter/backend/services"
)

// Define service objects at the package level
var (
	dataService *services.DataService
	startTime   time.Time
)

// InitServices initializes the service objects
func InitServices(database *db.Database) {
	dataService = services.NewDataService(database)
	startTime = time.Now()
	log.Println("API services initialized successfully")
}

// SetupRoutes configures the API routes
func SetupRoutes(r *gin.Engine) {
	// Health check endpoint - register at multiple paths for redundancy
	r.GET("/health", HealthCheck)

	log.Println("Health check endpoints registered at /health and /")

	// API routes
	api := r.Group("/api")
	{
		api.GET("/destinations/random", GetRandomDestination)
		api.POST("/users", CreateUser)
		api.GET("/users/:username", GetUser)

		// Game routes
		api.POST("/game/play", StartGame)
		api.GET("/game/:id/next-question", GetNextQuestion)
		api.POST("/game/:id/submit-answer", SubmitAnswer)
		api.GET("/game/:id/result", GetGameResult)
		api.GET("/game/:id/summary", GetGameSummary)
	}

	log.Println("All API routes registered successfully")
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	log.Println("Health check endpoint called")

	uptime := time.Since(startTime).String()

	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"uptime":    uptime,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// GetRandomDestination handles requests for a random destination
func GetRandomDestination(c *gin.Context) {
	destination, err := dataService.GetRandomDestination()
	if err != nil {
		fmt.Println("Error getting random destination:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random destination"})
		return
	}

	// Send the destination without modification
	c.JSON(http.StatusOK, destination)
}

// CreateUser handles requests to create a new user
func CreateUser(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := dataService.CreateUser(request.Username)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles requests to get a user by username
func GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := dataService.GetUser(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// StartGame handles requests to start a new game
func StartGame(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	gameID, err := dataService.CreateGame(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create game: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"game_id": gameID})
}

// GetNextQuestion handles requests to get the next question in a game
func GetNextQuestion(c *gin.Context) {
	gameID := c.Param("id")

	// Convert gameID to int
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	question, err := dataService.GetNextQuestion(gameIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to get next question: %v", err)})
		return
	}

	// Check if there are more questions
	hasNext, err := dataService.HasNextQuestion(gameIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to check for more questions: %v", err)})
		return
	}

	// Get destination details for each option
	optionsDisplay := make(map[int]string)
	for _, destID := range question.OptionDestinationIDs {
		dest, err := dataService.GetDestinationByID(destID)
		if err != nil {
			continue // Skip if destination not found
		}
		optionsDisplay[destID] = fmt.Sprintf("%s, %s", dest.City, dest.Country)
	}

	// Create response
	response := models.NextQuestionResponse{
		GameID:         gameIDInt,
		QuestionID:     question.ID,
		Question:       question.Question,
		OptionsDisplay: optionsDisplay,
		HasNext:        hasNext,
	}

	c.JSON(http.StatusOK, response)
}

// SubmitAnswer handles requests to submit an answer for a question
func SubmitAnswer(c *gin.Context) {
	gameID := c.Param("id")

	// Convert gameID to int
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	var request models.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate that the game ID in the URL matches the one in the request
	if gameIDInt != request.GameID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game ID mismatch"})
		return
	}

	result, err := dataService.SubmitAnswer(request.GameID, request.QuestionID, request.SelectedDestination)
	if err != nil {
		// Check for specific error messages
		if err.Error() == "question already answered" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question already answered"})
			return
		} else if err.Error() == "selected destination is not in the list of options" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Selected destination is not in the list of options"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to submit answer: %v", err)})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetGameResult handles requests to get the result of a game
func GetGameResult(c *gin.Context) {
	gameID := c.Param("id")

	// Convert gameID to int
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	result, err := dataService.GetGameResult(gameIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to get game result: %v", err)})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetGameSummary handles requests to get a summary of a game
func GetGameSummary(c *gin.Context) {
	gameID := c.Param("id")
	println("gameID", gameID)
	// Convert gameID to int
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid game ID"})
		return
	}

	summary, err := dataService.GetGameSummary(gameIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to get game summary: %v", err)})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ServeChallengePage serves the challenge page with proper Open Graph meta tags
func ServeChallengePage(c *gin.Context) {
	username := c.Param("username")
	gameID := c.Param("gameID")

	log.Printf("Serving challenge page for username: %s, gameID: %s", username, gameID)

	// Get the path to index.html
	webappPath := "../webapp/build"
	if _, err := os.Stat("/app/webapp/build"); err == nil {
		webappPath = "/app/webapp/build"
	}

	indexPath := filepath.Join(webappPath, "index.html")

	// Read the index.html file
	content, err := ioutil.ReadFile(indexPath)
	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		c.File(indexPath) // Fallback to regular file serving
		return
	}

	// Convert to string for easier manipulation
	htmlContent := string(content)

	// Define the image URL - use a reliable, consistent image
	imageURL := "https://images.unsplash.com/photo-1530521954074-e64f6810b32d?ixlib=rb-4.0.3&auto=format&fit=crop&w=1200&q=80"

	// Check if we already have OG tags
	if !strings.Contains(htmlContent, "og:image") {
		// Insert Open Graph meta tags before the closing head tag
		metaTags := fmt.Sprintf(`
    <meta property="og:title" content="Globetrotter Challenge" />
    <meta property="og:description" content="%s has challenged you to beat their score in Globetrotter!" />
    <meta property="og:image" content="%s" />
    <meta property="og:url" content="%s" />
    <meta property="og:type" content="website" />
    <meta property="og:site_name" content="Globetrotter" />
    <meta name="twitter:card" content="summary_large_image" />
  </head>`, username, imageURL, c.Request.URL.String())

		// Replace the closing head tag with our meta tags
		htmlContent = strings.Replace(htmlContent, "</head>", metaTags, 1)
	} else {
		// Update existing OG tags
		htmlContent = strings.Replace(htmlContent, `content="Globetrotter - The Ultimate Travel Guessing Game"`,
			fmt.Sprintf(`content="%s has challenged you to beat their score in Globetrotter!"`, username), 1)

		// If there's an existing og:image tag, update it
		if strings.Contains(htmlContent, `property="og:image"`) {
			startIdx := strings.Index(htmlContent, `property="og:image"`)
			endIdx := strings.Index(htmlContent[startIdx:], "/>") + startIdx
			if endIdx > startIdx {
				beforeTag := htmlContent[:startIdx]
				afterTag := htmlContent[endIdx+2:]
				htmlContent = beforeTag + fmt.Sprintf(`property="og:image" content="%s" />`, imageURL) + afterTag
			}
		}
	}

	// Serve the modified HTML content
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlContent)
}
