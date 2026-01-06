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
