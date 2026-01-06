This is a major pivot point. Up until now, we've been building a static website—it just displays fixed information.
Now, we are going to start turning it into a dynamic web application.
The core purpose of c500-web-go is e-commerce. The Buyer's Manual talks about browsing items and a checkout page, but those don't exist yet. The next logical step is to build the foundation for the online store.
To display products, our Go program first needs to understand what a "product" actually is. We need to define a data structure.
Here is the plan to start building the store:
 * Create a models directory: This will hold Go structs that define our data structures.
 * Define the Product Model: Create a Go struct representing an item for sale (Name, Price, ID, etc.).
 * Create Mock Data: Since we don't have a real database connected yet, we'll create a hardcoded list of products to test with.
Step 1 & 2: Create directory and define the Product Model
Create a new directory named models in your root folder.
Then, create a file inside it named product.go (c500-web-go/models/product.go).
This file defines the shape of our data. Notice the json:"..." tags—these will be useful later if we build a JSON API for the Discord bot.
package models

import "fmt"

// Product represents an item listed for sale.
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`      // Using float64 for simplicity for now
	ImageURL    string  `json:"image_url"`  // Path to product image
	SellerID    string  `json:"seller_id"`  // Discord ID of seller
}

// FormattedPrice returns a nicely formatted string of the price (e.g., "$10.00")
// This is a helper method we can use directly in our HTML templates.
func (p Product) FormattedPrice() string {
	return fmt.Sprintf("$%.2f", p.Price)
}

Step 3: Create Mock Data
Now we need some dummy data to display. In a real app, this would come from PostgreSQL, but for development, hardcoded data is fine.
We will add a simple function to models/product.go that returns a slice list of products.
Update c500-web-go/models/product.go to include this new function at the bottom:
package models

import "fmt"

// Product represents an item listed for sale.
type Product struct {
    // ... (same struct as above) ...
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	SellerID    string  `json:"seller_id"`
}

// FormattedPrice helper method
func (p Product) FormattedPrice() string {
	return fmt.Sprintf("$%.2f", p.Price)
}

// --- Mock Data Section ---

// GetMockProducts returns a hardcoded slice of products for testing purposes.
// In the future, this would be replaced by a database query.
func GetMockProducts() []Product {
	return []Product{
		{
			ID:          1,
			Name:        "Vintage Discord Wumpus Plush",
			Description: "A rare, collectible plushie from 2018. Slightly worn tag.",
			Price:       45.00,
			ImageURL:    "https://placehold.co/400x300/png?text=Wumpus+Plush", // Using a placeholder image service
			SellerID:    "discord_user_123",
		},
		{
			ID:          2,
			Name:        "Custom Server Emoji Pack (x10)",
			Description: "I will design 10 custom, high-quality emojis for your server. Digital delivery.",
			Price:       15.50,
			ImageURL:    "https://placehold.co/400x300/png?text=Emoji+Pack",
			SellerID:    "artist_jane#5555",
		},
		{
			ID:          3,
			Name:        "Level 3 Server Boost Build",
			Description: "Consultation and setup to help get your community server to Level 3 status.",
			Price:       99.99,
			ImageURL:    "https://placehold.co/400x300/png?text=Server+Boost",
			SellerID:    "admin_mike",
		},
	}
}

Now our application knows what a product looks like and has some data to work with. The next logical step will be to create the handlers and templates to actually display this data on a "Store" page.
