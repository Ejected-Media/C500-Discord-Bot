This is a major milestone. You have built a complete, functional prototype of the web frontend and the payment integration.
Right now, however, your application has a significant limitation: it has no memory.
If you restart the server, any changes (if we had features to make them) would reset. The products are just hardcoded text in a Go file.
The logical next step is to build the foundation for a real application by replacing the mock data with a real, persistent database. We will use PostgreSQL, a robust, open-source SQL database that is perfect for e-commerce data.
Here is the plan to initialize the database layer:
 * Prerequisite: Install PostgreSQL.
 * Database Setup: Create the actual database and define the table structure (schema) for products.
 * Environment Config: Update .env with database credentials.
 * Install Go Driver: Get the package that allows Go to talk to PostgreSQL.
 * Connect in main.go: Test the connection on server startup.
Step 1: Prerequisite - Install PostgreSQL
If you don't have PostgreSQL installed on your development machine yet, you need to do that first.
 * Mac (via Homebrew): brew install postgresql@14 then brew services start postgresql@14`
 * Windows/Linux: Download the installer from the official PostgreSQL downloads page.
Tip: I highly recommend installing a GUI tool like pgAdmin 4 (often comes with the installer) or TablePlus to easily manage your database visually.
Step 2: Database Setup (Create Schema)
We need to tell PostgreSQL what our data looks like. We'll create a SQL script to define our database and the products table.
 * Create a new directory in your project root named db.
 * Inside it, create a file named schema.sql (c500-web-go/db/schema.sql).
 * Paste the following SQL commands into it:
<!-- end list -->
-- db/schema.sql

-- 1. Create the database (you might need to run this line separately depending on your setup)
-- CREATE DATABASE c500_ecommerce;

-- Connect to the 'c500_ecommerce' database before running the rest.

-- 2. Create the Products table
-- We use SERIAL for the ID to make it auto-incrementing.
-- We use DECIMAL(10, 2) for price to handle currency accurately (never use floats for money in DBs!)
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    image_url VARCHAR(512),
    seller_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Seed some initial test data (so the store isn't empty)
INSERT INTO products (name, description, price, image_url, seller_id) VALUES
('DB Vintage Wumpus Plush', 'From the database! A rare 2018 plushie.', 45.00, 'https://placehold.co/600x400/png?text=DB+Wumpus', 'discord_user_123'),
('DB Emoji Pack', '10 custom emojis stored in SQL.', 15.50, 'https://placehold.co/600x400/png?text=DB+Emoji', 'artist_jane'),
('DB Server Boost', 'Level 3 setup, fetched dynamically.', 99.99, 'https://placehold.co/600x400/png?text=DB+Boost', 'admin_mike');

Action: Run these SQL commands against your local PostgreSQL installation to create the database and tables.
Step 3: Environment Config (.env)
We need to tell our Go program how to find the database. Open your .env file and add your connection details.
Note: Adjust username, password, and host if you configured them differently during installation.
# c500-web-go/.env

# ... existing Stripe keys ...
STRIPE_PUBLISHABLE_KEY=...
STRIPE_SECRET_KEY=...
DOMAIN=http://localhost:8080

# Database Connection Info
# Format: postgres://username:password@host:port/database_name?sslmode=disable
# Example for default local setup:
DATABASE_URL=postgres://postgres:password@localhost:5432/c500_ecommerce?sslmode=disable

Step 4: Install Go Database Driver
Go's standard library has database interfaces, but it needs a specific driver to talk to PostgreSQL. We'll use the reliable lib/pq driver.
Stop your server and run:
go get github.com/lib/pq

Step 5: Connect in main.go
Now, let's modify main.go to open a connection pool to the database when the server starts. We won't use it to fetch data yet, but we will verify that the connection works.
Open main.go and apply these changes:
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

Test the Connection
Run your server:
go run main.go

Watch the terminal output closely. You are looking for this success message right before the server starts:
Successfully connected to PostgreSQL database!
Starting server on http://localhost:8080 ...

If you see that message, congratulations! Your Go application is now connected to a real, persistent SQL database. The next step will be refactoring your models to stop using mock data and start running real SQL queries.
