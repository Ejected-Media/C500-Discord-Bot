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
