package main

import (
    // Import the database/sql package to define interfaces
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"c500-web-go/handlers"

	"github.com/joho/godotenv"
    // Import the driver anonymously (using _) so it registers itself
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file.")
	}

	// --- NEW SECTION: Database Connection ---
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set.")
	}

	// Open establishes a connection pool structure, but doesn't connect yet.
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// Close the connection when main() exits (e.g., server shutdown)
	defer db.Close()

	// Ping actually tries to connect to verify credentials are correct.
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v. Check your .env credentials.", err)
	}
	log.Println("Successfully connected to PostgreSQL database!")
	// ----------------------------------------


	// 2. Parse Templates (Same as before)
    // ... rest of main.go is unchanged for now ...
  
