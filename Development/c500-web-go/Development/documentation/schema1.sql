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
