This is the moment where your application truly comes alive.
Right now, your app is connected to the database, but it's ignoring it. It's still using the hardcoded mock data from step 11.
The next step is to refactor the application to fetch data from PostgreSQL.
We need to update our models package to interact with the database (*sql.DB) and then update our handlers to use these new database-aware methods.
Here is the plan to switch from mock data to real data:
 * Refactor models/product.go: Replace mock data with a struct that holds the database connection and methods to execute SQL queries.
 * Inject Dependencies in main.go: Initialize the new product model with the database connection and pass it to the handlers.
 * Update handlers/store.go: Use the new model to fetch all products for the main store page.
 * Update handlers/checkout.go: Use the new model to fetch a specific product for checkout.
Step 1: Refactor Product Model (models/product.go)
Open c500-web-go/models/product.go. We are going to delete all the mock data and replace it with real database logic.
We will define a ProductModel struct that holds the database connection pool. This allows us to define methods on it that can run queries.
package models

import (
	"database/sql"
	"fmt"
)

// Product represents an item listed for sale in the DB.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	// Note: scanning SQL DECIMAL into float64 is okay for display, 
    // but for complex financial math, consider a dedicated decimal package later.
	Price       float64 `json:"price"` 
	ImageURL    string  `json:"image_url"`
	SellerID    string  `json:"seller_id"`
}

// FormattedPrice helper method (remains the same)
func (p Product) FormattedPrice() string {
	return fmt.Sprintf("$%.2f", p.Price)
}

// --- NEW DB LOGIC ---

// ProductModel holds the database connection pool.
// This is how we will interact with the 'products' table.
type ProductModel struct {
	DB *sql.DB
}

// All fetches all products from the database.
func (m *ProductModel) All() ([]Product, error) {
	// 1. Write the SQL query
	stmt := `SELECT id, name, description, price, image_url, seller_id FROM products`

	// 2. Execute the query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Ensure the connection is released back to the pool when we are done
	defer rows.Close()

	// 3. Iterate through the results
	var products []Product
	for rows.Next() {
		var p Product
		// Scan copies the columns from the current row into the struct fields
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SellerID)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// Check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Get fetches a single product by its ID.
func (m *ProductModel) Get(id int) (*Product, error) {
	// 1. Write query with a placeholder ($1) to prevent SQL injection
	stmt := `SELECT id, name, description, price, image_url, seller_id FROM products WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)

	var p Product
	// 2. Scan the single result row
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SellerID)
	if err != nil {
		// This error occurs specifically if no rows were found
		if err == sql.ErrNoRows {
			return nil, nil // Return nil product and nil error (not found, but not a system failure)
		}
		return nil, err
	}

	return &p, nil
}

Step 2: Inject Dependencies in main.go
Now we need to wire everything up. In main.go, after connecting to the DB, we create an instance of our new ProductModel and pass it to the handlers that need data access.
Open c500-web-go/main.go:
// ... imports remain the same ...
import (
    // ...
	"c500-web-go/handlers"
	"c500-web-go/models" // Ensure models is imported

	_ "github.com/lib/pq"
)

func main() {
    // 1. Load Env (same) ...
    // ...

    // --- DB Connection (same as last step) ---
	dbURL := os.Getenv("DATABASE_URL")
    // ... checks ...
	db, err := sql.Open("postgres", dbURL)
	if err != nil { log.Fatalf("%v", err) }
	defer db.Close()

	err = db.Ping()
	if err != nil { log.Fatalf("%v", err) }
	log.Println("Successfully connected to PostgreSQL database!")
    // ----------------------------------------

    // --- NEW: Initialize Data Models ---
    // Create the model instance holding the DB connection pool
    productModel := &models.ProductModel{DB: db}


	// 2. Parse Templates (same) ...
    tmpl, err := parseTemplates("./templates")
    // ...


	// 3. Initialize Handlers
    // We need to update NewStoreHandler and NewCheckoutHandler to accept the productModel.
    // See steps 3 & 4 below for the changes to those files.

	homeHandler := handlers.NewHomeHandler(tmpl)
	docsHandler := handlers.NewDocsHandler(tmpl)
    // CHANGE: Pass model to store handler
	storeHandler := handlers.NewStoreHandler(tmpl, productModel) 
    // CHANGE: Pass model to checkout handler
	checkoutHandler := handlers.NewCheckoutHandler(productModel) 
	orderHandler := handlers.NewOrderHandler(tmpl)

    // ... rest of main.go remains the same ...
}

Step 3: Update Store Handler (handlers/store.go)
Update store.go to accept the ProductModel and use its All() method.
package handlers

import (
	"html/template"
	"log"
	"net/http"
    "strconv"

	"c500-web-go/models"
)

// StoreHandler holds dependencies.
type StoreHandler struct {
	Templates *template.Template
    // NEW: Add the product model dependency
	Products  *models.ProductModel
}

// Update constructor to accept the model
func NewStoreHandler(t *template.Template, p *models.ProductModel) *StoreHandler {
	return &StoreHandler{
		Templates: t,
		Products:  p,
	}
}

// ... StorePageData struct remains the same ...

// Index handles the main store page.
func (h *StoreHandler) Index(w http.ResponseWriter, r *http.Request) {
	// --- CHANGE: Fetch data from real DB instead of mock functions ---
	products, err := h.Products.All()
	if err != nil {
		log.Printf("Error fetching products from DB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := StorePageData{
		Title:    "Browse Items - c500-web-go Store",
		Products: products,
	}
    // ... rendering logic remains the same ...
	h.render(w, "store.html", data) // Assuming you refactored render helper, otherwise use ExecuteTemplate
}

// ProductDetail handles single product page.
func (h *StoreHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

    // --- CHANGE: Use DB model to get product ---
	product, err := h.Products.Get(id)
	if err != nil {
        // Log DB errors
		log.Printf("Error fetching product %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
    // If product is nil, it means it wasn't found in DB
	if product == nil {
		http.NotFound(w, r)
		return
	}

    // ... rendering logic remains the same, passing *product ...
     h.render(w, "product_detail.html", product)
}

// (Optional: If you haven't refactored the 'render' helper from docs.go into a shared place, 
// you'll need to copy/paste it here or use ExecuteTemplate directly like before. 
// For brevity, I assumed a helper function like in docs.go)
func (h *StoreHandler) render(w http.ResponseWriter, tmplName string, data interface{}) {
     // ... (implementation same as in docs.go but taking interface{}) ...
	if h.Templates.Lookup(tmplName) == nil {
		log.Printf("Error: Template '%s' not found.", tmplName)
		http.Error(w, "Template missing", http.StatusInternalServerError)
		return
	}
	err := h.Templates.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("Error executing template '%s': %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

Step 4: Update Checkout Handler (handlers/checkout.go)
Update checkout.go to get the real product data before sending it to Stripe.
package handlers

import (
	// ... imports ...
	"c500-web-go/models"
    // ...
)

// CheckoutHandler dependencies.
type CheckoutHandler struct{
    // NEW: Add model dependency
    Products *models.ProductModel
}

// Update constructor
func NewCheckoutHandler(p *models.ProductModel) *CheckoutHandler {
	return &CheckoutHandler{
        Products: p,
    }
}

func (h *CheckoutHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
    // ... get ID from URL (same) ...
	idStr := r.PathValue("id")
    id, _ := strconv.Atoi(idStr) // (simplified error check for brevity here)

	// 2. Fetch real product details from DB
	product, err := h.Products.Get(id)
	if err != nil {
        log.Printf("DB Error: %v", err)
		http.Error(w, "Internal Error", 500)
		return
	}
    if product == nil {
        http.NotFound(w, r)
        return
    }

	// ... Stripe configuration (same) ...

	// 4. Create session parameters using REAL product data
	params := &stripe.CheckoutSessionParams{
		// ... types ...
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        // Use the data from the DB product
						Name:        stripe.String(product.Name),
						Description: stripe.String(product.Description),
						Images:      stripe.StringSlice([]string{product.ImageURL}),
					},
                    // Use the real price from DB
					UnitAmount: stripe.Int64(int64(product.Price * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
        // ... rest of Stripe params (same) ...
	}
    
    // ... create session and redirect (same) ...
}

Test It Out
 * Restart your server: go run main.go
 * Visit the store: http://localhost:8080/store
You should see the three products you inserted into your database via the schema.sql script. The data is no longer coming from Go code; it's coming from PostgreSQL. You have a real, database-driven web application!
