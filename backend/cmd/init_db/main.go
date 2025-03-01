package main

import (
	"log"
	"os"

	"github.com/shubhsherl/globetrotter/backend/db"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/globetrotter.db"
	}
	
	log.Printf("Initializing database at %s...", dbPath)
	
	// Remove existing database if it exists
	if _, err := os.Stat(dbPath); err == nil {
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Failed to remove existing database: %v", err)
		}
		log.Println("Removed existing database")
	}
	
	// Initialize database
	if err := db.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	log.Println("Database initialized successfully")
} 