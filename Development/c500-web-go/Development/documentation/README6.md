This is a pivotal moment. We have successfully transitioned from a static documentation site to a dynamic web application that can render data.
Currently, you have a store page that lists all products. The natural next step in any e-commerce flow is to allow the user to click on an individual product to see more details.
We are going to build the Product Detail Page.
This requires a significant upgrade to our application's capabilities: we need to handle dynamic URLs. Instead of just handling /store, we need to handle /store/product/1, /store/product/2, etc., where the number at the end changes based on which product the user clicked.
Fortunately, Go 1.22 introduced powerful new routing capabilities in the standard library that make this easy.
Here is the plan:
 * Update Model: Add a function to find a single product by its ID in our mock data.
 * Update Handler: Create a new handler function that extracts the ID from the URL URL, fetches the correct product, and renders a template.
 * Create Template: Design the dedicated product detail HTML page.
 * Update Main & Links: Register the new dynamic route and connect the buttons on the main store page.
Step 1: Update Product Model (models/product.go)
We need a way to retrieve just one product. Open c500-web-go/models/product.go and add the GetProductByID function to the bottom of the file.
package models

import (
	"errors"
	"fmt"
)

// ... (Product struct and FormattedPrice method remain the same) ...

// --- Mock Data Section ---

var mockProducts = []Product{
	{
		ID:          1,
		Name:        "Vintage Discord Wumpus Plush",
		Description: "A rare, collectible plushie from 2018. Slightly worn tag. Ships immediately.",
		Price:       45.00,
		ImageURL:    "https://placehold.co/600x400/png?text=Wumpus+Plush",
		SellerID:    "discord_user_123",
	},
	{
		ID:          2,
		Name:        "Custom Server Emoji Pack (x10)",
		Description: "I will design 10 custom, high-quality emojis for your server. Digital delivery within 48 hours.",
		Price:       15.50,
		ImageURL:    "https://placehold.co/600x400/png?text=Emoji+Pack",
		SellerID:    "artist_jane#5555",
	},
	{
		ID:          3,
		Name:        "Level 3 Server Boost Build",
		Description: "Complete consultation and community setup to help get your server to Level 3 status. Includes bot configuration.",
		Price:       99.99,
		ImageURL:    "https://placehold.co/600x400/png?text=Server+Boost",
		SellerID:    "admin_mike",
	},
}

// GetMockProducts returns all products.
func GetMockProducts() []Product {
	return mockProducts
}

// --- NEW FUNCTION ---

// GetProductByID searches the mock data for a product matching the given ID.
// It returns the Product and nil error if found, or an empty Product and an error if not.
func GetProductByID(id int) (Product, error) {
	for _, p := range mockProducts {
		if p.ID == id {
			return p, nil
		}
	}
	// Return an empty product and a standardized error if ID doesn't exist.
	return Product{}, errors.New("product not found")
}

Step 2: Update Store Handler (handlers/store.go)
Now we need a handler that can read the URL path. In Go 1.22+, we can use r.PathValue("id") to grab dynamic parts of the URL.
Open c500-web-go/handlers/store.go and add the ProductDetail method.
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv" // Need strconv to convert string ID from URL to int

	"c500-web-go/models"
)

// ... (StoreHandler struct, NewStoreHandler, StorePageData, and Index method remain the same) ...

// --- NEW HANDLER METHOD ---

// ProductDetail handles requests for a single product page (e.g., /store/product/1).
func (h *StoreHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the ID from the URL path.
	// This relies on the route being defined as "/store/product/{id}" in main.go.
	idStr := r.PathValue("id")

	// 2. Convert string ID to integer.
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r) // Return 404 if ID is not a valid positive integer
		return
	}

	// 3. Fetch the specific product using our new model function.
	product, err := models.GetProductByID(id)
	if err != nil {
		// If product isn't found in our "database", return 404.
		http.NotFound(w, r)
		return
	}

	// 4. Render the template specifically for detail view.
	// We pass the product directly as the data dot (.).
	tmplName := "product_detail.html"
	if h.Templates.Lookup(tmplName) == nil {
		log.Printf("Error: Template '%s' not found.", tmplName)
		http.Error(w, "Template missing", http.StatusInternalServerError)
		return
	}

	err = h.Templates.ExecuteTemplate(w, tmplName, product)
	if err != nil {
		log.Printf("Error executing template '%s': %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

Step 3: Create Product Detail Template (templates/product_detail.html)
Create this new HTML file in your templates/ directory. It's similar to the store card, but designed to showcase a single item prominently.
{{extends "base.html"}}

{{define "title"}}{{.Name}} - c500-web-go Store{{end}}

{{define "head_extra"}}
<style>
    .product-detail-container {
        display: grid;
        grid-template-columns: 1fr 1fr; /* Two equal columns */
        gap: 40px;
        margin-top: 2rem;
    }

    /* Make grid single column on smaller screens */
    @media (max-width: 768px) {
        .product-detail-container { grid-template-columns: 1fr; }
    }

    .detail-image img {
        width: 100%;
        border-radius: 8px;
        border: 1px solid var(--color-border);
    }

    .detail-info h1 {
        font-size: 2.5rem;
        margin-bottom: 1rem;
        color: var(--color-header-bg);
    }

    .detail-price {
        font-size: 2rem;
        color: var(--color-primary);
        font-weight: bold;
        margin-bottom: 1.5rem;
    }

    .detail-description {
        font-size: 1.1rem;
        line-height: 1.8;
        margin-bottom: 2rem;
        color: #555;
    }

    .seller-info {
        margin-bottom: 2rem;
        font-style: italic;
        color: #777;
    }

    /* Re-using the buy-button style, but making it bigger */
    .buy-now-btn {
        display: inline-block;
        padding: 15px 40px;
        font-size: 1.2rem;
        background-color: var(--color-primary);
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        transition: background-color 0.2s;
    }
    .buy-now-btn:hover { background-color: #357abd; text-decoration: none; }
    
    .back-link { display: block; margin-bottom: 1rem; }
</style>
{{end}}

{{define "content"}}
<a href="/store" class="back-link">&larr; Back to Store</a>

<article class="product-detail-container">
    <div class="detail-image">
        <img src="{{.ImageURL}}" alt="{{.Name}}">
    </div>

    <div class="detail-info">
        <h1>{{.Name}}</h1>
        <p class="detail-price">{{.FormattedPrice}}</p>
        
        <div class="detail-description">
            <p>{{.Description}}</p>
        </div>

        <p class="seller-info">Sold by: {{.SellerID}}</p>

        <a href="#" class="buy-now-btn">Buy Now</a>
    </div>
</article>
{{end}}

Step 4: Update Main & Links
We have two final small tasks to connect everything.
1. Update main.go to register the dynamic route.
Crucial: Notice the syntax "GET /store/product/{id}". The {id} tells Go that this part of the URL is dynamic, and the GET ensures this handler only responds to GET requests.
// ... inside main.go ...

	// --- Application Routes ---
	mux.HandleFunc("/", homeHandler.Landing)

	// Main Store Listings
	mux.HandleFunc("/store", storeHandler.Index)

	// --- NEW ROUTE: Product Detail ---
    // The {id} is a path value wildcard introduced in Go 1.22
	mux.HandleFunc("GET /store/product/{id}", storeHandler.ProductDetail)

	// ... docs routes remain the same ...

2. Update templates/store.html to link the buttons.
Open templates/store.html and change the placeholder <button> into a real <a> tag that links to the product's unique URL.
<div class="product-details">
    <h2>{{.Name}}</h2>
    <p class="price">{{.FormattedPrice}}</p>
    <p class="description">{{.Description}}</p>
    <p class="seller">Seller: {{.SellerID}}</p>
    
    <a href="/store/product/{{.ID}}" class="buy-button" style="text-align: center; text-decoration: none; display: block;">View Details</a>
</div>

Test It Out
Restart your server: go run main.go
 * Go to http://localhost:8080/store.
 * Click the "View Details" button on the "Vintage Discord Wumpus Plush".
 * You should be taken to http://localhost:8080/store/product/1 and see the detailed view of that specific item.
 * Try changing the URL manually to /store/product/99. You should correctly see a "404 Not Found" error page provided by the browser.
 
