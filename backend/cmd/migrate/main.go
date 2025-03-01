package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbPath        = "./data/globetrotter.db"
	migrationsDir = "./migrations"
)

func main() {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get list of applied migrations
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		log.Fatalf("Failed to query migrations: %v", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("Failed to scan migration row: %v", err)
		}
		appliedMigrations[name] = true
	}

	// Get list of migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Sort migration files by name
	sort.Strings(migrationFiles)

	// Apply migrations
	for _, file := range migrationFiles {
		if appliedMigrations[file] {
			fmt.Printf("Migration %s already applied\n", file)
			continue
		}

		fmt.Printf("Applying migration %s...\n", file)

		// Read migration file
		content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		// Execute migration
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}

		// Record migration
		_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", file)
		if err != nil {
			log.Fatalf("Failed to record migration %s: %v", file, err)
		}

		fmt.Printf("Migration %s applied successfully\n", file)
	}

	fmt.Println("All migrations applied successfully")
}
