package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// IMPORT the new package
	"github.com/joho/godotenv"

	"c500-web-go/handlers"
)

func main() {
	// --- NEW SECTION: Load Environment Variables ---
	// Load values from .env file into the system environment
	err := godotenv.Load()
	if err != nil {
		// It's okay if .env doesn't exist in production (we'd set real env vars there),
		// but for local dev, we want to know if it failed.
		log.Println("Warning: Error loading .env file. Ensure STRIPE_SECRET_KEY is set in environment.")
	}

	// Verify keys are loaded (DO NOT PRINT THE ACTUAL KEYS IN PRODUCTION LOGS)
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		log.Fatal("Error: STRIPE_SECRET_KEY environment variable is not set.")
	}
	log.Println("Environment variables loaded successfully.")
	// -----------------------------------------------


	// 1. Parse Templates (remains the same)
	tmpl, err := parseTemplates("./templates")
    // ... rest of the file remains the same ...
  
