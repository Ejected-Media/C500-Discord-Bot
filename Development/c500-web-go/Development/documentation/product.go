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
