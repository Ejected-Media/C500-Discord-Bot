This is a critical turning point. You have successfully established a secure communication channel between Stripe and your local server via webhooks.
Right now, your server just logs the payment and forgets it. A real e-commerce application needs a permanent record of every sale.
The next step is to create an "Orders" database table and save successful transactions into it.
By doing this, you will officially move from a stateless prototype to a stateful application that tracks business history.
Here is the plan to record sales:
 * Update Database Schema: Create a new orders table in PostgreSQL to store transaction details.
 * Create Order Model: Write Go code to define the Order data structure and a function to insert a new order into the database.
 * Update Checkout Handler: Pass the internal Product ID to Stripe as "metadata" so Stripe can send it back to us later.
 * Update Webhook Handler: Modify the webhook to read that metadata and save the order to our new database table.
 * Update main.go: Wire up the new Order model dependency.
Step 1: Update Database Schema (db/schema.sql)
We need a place to store orders. Add the following SQL to your existing db/schema.sql file (or run it directly in your database tool).
This table links back to the products table using a Foreign Key, ensuring data integrity.
-- db/schema.sql

-- ... existing products table creation ...

-- 4. Create the Orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    stripe_session_id VARCHAR(255) UNIQUE NOT NULL, -- The unique ID from Stripe for this transaction
    product_id INT NOT NULL REFERENCES products(id), -- Links back to the product bought
    amount_total INT NOT NULL, -- Stored in cents (e.g., 4500 for $45.00)
    customer_email VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL, -- e.g., 'paid'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

Action: Run this new CREATE TABLE command against your PostgreSQL database to create the table.
Step 2: Create Order Model (models/order.go)
Create a new file named c500-web-go/models/order.go.
This file defines what an Order looks like in Go and provides the method to save it to the database.
package models

import (
	"database/sql"
	"time"
)

// Order represents a completed transaction.
type Order struct {
	ID              int       `json:"id"`
	StripeSessionID string    `json:"stripe_session_id"`
	ProductID       int       `json:"product_id"`
	AmountTotal     int       `json:"amount_total"` // in cents
	CustomerEmail   string    `json:"customer_email"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// OrderModel holds the database connection pool.
type OrderModel struct {
	DB *sql.DB
}

// Insert adds a new order record to the database.
func (m *OrderModel) Insert(order *Order) error {
	stmt := `
		INSERT INTO orders (stripe_session_id, product_id, amount_total, customer_email, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	// Use QueryRow because we expect one row back (the new ID and timestamp)
	err := m.DB.QueryRow(
		stmt,
		order.StripeSessionID,
		order.ProductID,
		order.AmountTotal,
		order.CustomerEmail,
		order.Status,
	).Scan(&order.ID, &order.CreatedAt) // Update the struct with the new DB ID and time

	if err != nil {
		return err
	}
	return nil
}

Step 3: Update Checkout Handler (handlers/checkout.go)
This is a crucial step. When we create the checkout session, we need to tell Stripe which internal ID (product.ID) is being bought. Stripe allows us to attach arbitrary key-value pairs called metadata to a session. They will send this data back to us in the webhook.
Open handlers/checkout.go and add the Metadata field to the CheckoutSessionParams:
// ... inside handlers/checkout.go ...
	// ... imports, including strconv ...
    "strconv"
    // ...

func (h *CheckoutHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
    // ... fetch product from DB ...
	product, err := h.Products.Get(id)
    // ... checks ...

	// ... stripe setup ...

	// 4. Create the checkout session parameters
	params := &stripe.CheckoutSessionParams{
		// --- NEW SECTION: Attach Metadata ---
		// We attach our internal Product ID so we know what was bought later.
		Metadata: map[string]string{
			"product_id": strconv.Itoa(product.ID),
		},
		// ------------------------------------

		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
        // ... rest of the params remain the same ...

Step 4: Update Webhook Handler (handlers/webhook.go)
Now update the webhook handler to read that metadata and use the new OrderModel to save the data.
Open handlers/webhook.go:
package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv" // Needed for converting metadata string to int

	"c500-web-go/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// WebhookHandler dependencies updated to include Orders.
type WebhookHandler struct {
	Products *models.ProductModel
	Orders   *models.OrderModel // <-- Add this
}

// Update constructor to accept OrderModel.
func NewWebhookHandler(p *models.ProductModel, o *models.OrderModel) *WebhookHandler {
	return &WebhookHandler{
		Products: p,
		Orders:   o, // <-- Store it
	}
}

// Handle endpoint... (beginning of function remains the same) ...
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
    // ... read body, verify signature, parse event JSON ...
    // ... down to the event type check:

	// 4. Handle the specific event type
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// --- BUSINESS LOGIC UPDATED ---

		// 1. Extract the Product ID from the Metadata we sent earlier
		productIDStr := session.Metadata["product_id"]
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			log.Printf("Error: Malformed product_id in metadata: %v", err)
			// We still return 200 OK to Stripe so they don't keep retrying,
			// but we log the error for our own investigation.
			w.WriteHeader(http.StatusOK)
			return
		}

		// 2. Create the Order struct payload
		order := &models.Order{
			StripeSessionID: session.ID,
			ProductID:       productID,
			AmountTotal:     int(session.AmountTotal), // Stripe sends int64, cast to int for DB
			CustomerEmail:   session.CustomerDetails.Email,
			Status:          string(session.PaymentStatus), // e.g., "paid"
		}

		// 3. Insert into PostgreSQL
		err = h.Orders.Insert(order)
		if err != nil {
			log.Printf("DATABASE ERROR: Could not save order: %v", err)
			// In a real app, you might return 500 here to make Stripe retry later,
			// depending on your idempotency strategy.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Order successfully saved to DB! ID: %d, Email: %s", order.ID, order.CustomerEmail)
	}

	// 5. Acknowledge receipt
	w.WriteHeader(http.StatusOK)
}

Step 5: Update main.go to Wire Dependencies
Finally, update main.go to initialize the new OrderModel and pass it to the webhook handler.
// ... imports and main() setup ...

	// ... after connecting to DB ...
	log.Println("Successfully connected to PostgreSQL database!")

	// --- Initialize Data Models ---
	productModel := &models.ProductModel{DB: db}
    // NEW: Init OrderModel
	orderModel := &models.OrderModel{DB: db}

	// ... parse templates ...

	// --- Initialize Handlers ---
    // ... other handlers ...
	orderHandler := handlers.NewOrderHandler(tmpl)

	// CHANGE: Update Webhook handler constructor to pass orderModel
	webhookHandler := handlers.NewWebhookHandler(productModel, orderModel)

	// ... router setup and server start remain the same ...

Final Test Run
 * Ensure Stripe CLI is running: stripe listen ... should still be active in a terminal.
 * Restart Go Server: go run main.go
 * Perform a Test Purchase in your browser.
 * Check the Go Server Terminal: You should see a new success message:
   > ✅ Order successfully saved to DB! ID: 1, Email: test@example.com
   > 
 * Verify in Database: Open your database tool (pgAdmin, TablePlus, or psql command line) and run:
   SELECT * FROM orders;

   You will see the record of the purchase you just made, linked to the correct product ID.
You have now built a complete, end-to-end e-commerce flow that persists data!
