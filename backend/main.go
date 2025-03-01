package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shubhsherl/globetrotter/backend/api"
	"github.com/shubhsherl/globetrotter/backend/db"
)

func main() {
	// Set up database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/globetrotter.db"
	}

	if err := db.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize database instance
	database := db.GetDB()

	// Initialize API services
	api.InitServices(database)

	// Set up Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Setup routes
	api.SetupRoutes(r)

	// Serve static files for production
	r.Static("/static", "../webapp/build/static")
	r.StaticFile("/", "../webapp/build/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("../webapp/build/index.html")
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
