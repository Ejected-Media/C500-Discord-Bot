This is the most exciting step yet. We are going to turn this into a real e-commerce application by integrating Stripe payments.
Currently, the "Buy Now" button on your product detail page does nothing. Our goal is to make it so that when clicked, the user is redirected to a secure Stripe-hosted checkout page to pay for that specific item.
This is a multi-step process. We will start with the foundational setup: getting your API keys, securing them, and installing the necessary Go libraries.
Here is the plan for this phase:
 * Prerequisite: Get Stripe Test API Keys.
 * Environment Setup: Create a .env file to store secrets securely.
 * Install Dependencies: Install the Stripe Go client and a library to read .env files.
 * Update main.go: Load the environment variables on startup.
Step 1: Prerequisite - Get Stripe Test API Keys
Before writing code, you need an account with Stripe.
 * Go to dashboard.stripe.com/register and sign up (it's free).
 * Once logged in, look at the top right corner of the dashboard and ensure the "Test mode" toggle is switched ON.
   * Crucial: Never use live keys for development.
 * Go to the Developers section, then API keys.
 * You need two keys:
   * Publishable key: Starts with pk_test_... (Safe to share in frontend code).
   * Secret key: Starts with sk_test_... (Keep this safe! Never commit it to Git).
Step 2: Environment Setup (.env)
It is a crucial security best practice never to hardcode API secrets directly in your Go code. Instead, we use environment variables.
 * In the root directory of your project (c500-web-go/), create a new file named exactly .env (dot env).
 * Open this file and paste your Stripe keys in the following format:
<!-- end list -->
# c500-web-go/.env

# Stripe API Keys (TEST MODE)
STRIPE_PUBLISHABLE_KEY=pk_test_YOUR_ACTUAL_PUBLISHABLE_KEY_HERE
STRIPE_SECRET_KEY=sk_test_YOUR_ACTUAL_SECRET_KEY_HERE

# The base URL of your local server (needed for Stripe redirects later)
DOMAIN=http://localhost:8080

Important: If you are using Git, add .env to your .gitignore file immediately so you don't accidentally upload your secrets.
Step 3: Install Go Dependencies
We need two new Go packages:
 * github.com/stripe/stripe-go/v76: The official Stripe client library for Go.
 * github.com/joho/godotenv: A handy library that reads a local .env file and loads the values into Go's environment variables.
Stop your server and run these commands in your terminal:
go get github.com/stripe/stripe-go/v76
go get github.com/joho/godotenv

Step 4: Update main.go to Load Secrets
Now we need to tell our application to load that .env file when it starts up, so the keys are available to our code.
Open main.go and add the following to the very beginning of the main() function.
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

Test It Out
Run your server:
go run main.go

Watch the terminal output. You should see:
Environment variables loaded successfully.
Templates parsed successfully.
Starting server on http://localhost:8080 ...

If you see the "Error: STRIPE_SECRET_KEY..." message, double-check that your .env file is in the correct root directory and named correctly.
We are now ready to write the code that actually creates a Stripe Checkout session!
