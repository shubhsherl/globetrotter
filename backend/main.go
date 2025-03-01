package main

import (
	"log"
	"os"
	"path/filepath"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shubhsherl/globetrotter/backend/api"
	"github.com/shubhsherl/globetrotter/backend/db"
)

func main() {
	log.Println("Starting Globetrotter application...")

	// Set up database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/globetrotter.db"
	}
	log.Printf("Using database at: %s", dbPath)

	// Ensure the database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	if err := db.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Initialize database instance
	database := db.GetDB()

	// Initialize API services
	api.InitServices(database)
	log.Println("API services initialized")

	// Set up Gin router
	if gin.Mode() == gin.ReleaseMode {
		log.Println("Running in release mode")
	} else {
		log.Println("Running in debug mode")
	}
	r := gin.Default()

	// Configure CORS - Commented out since frontend and backend are on the same domain
	/*
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type"},
			AllowCredentials: true,
		}))
	*/

	// Setup routes
	api.SetupRoutes(r)
	log.Println("API routes configured")

	// Check if we're running in Docker
	webappPath := "../webapp/build"
	if _, err := os.Stat("/app/webapp/build"); err == nil {
		webappPath = "/app/webapp/build"
		log.Println("Using Docker webapp path:", webappPath)
	}

	// Serve static files for production
	r.Static("/static", filepath.Join(webappPath, "static"))
	r.StaticFile("/favicon.ico", filepath.Join(webappPath, "favicon.ico"))
	r.StaticFile("/index.html", filepath.Join(webappPath, "index.html"))

	// NoRoute should come after explicit static file handlers
	r.NoRoute(func(c *gin.Context) {
		log.Printf("No route found for: %s, serving index.html", c.Request.URL.Path)
		c.File(filepath.Join(webappPath, "index.html"))
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
