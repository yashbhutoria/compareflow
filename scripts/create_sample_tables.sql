-- Create sample tables for testing validations
-- These simulate production and data warehouse tables

-- Create a sample products table (production)
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create a sample products table for data warehouse (with slight differences)
CREATE TABLE IF NOT EXISTS products_dw (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP
);

-- Insert sample data into products
INSERT INTO products (name, price, category) VALUES
    ('Laptop Pro 15"', 1299.99, 'Electronics'),
    ('Wireless Mouse', 29.99, 'Electronics'),
    ('Office Chair', 249.99, 'Furniture'),
    ('Standing Desk', 599.99, 'Furniture'),
    ('USB-C Cable', 19.99, 'Electronics'),
    ('Coffee Maker', 79.99, 'Appliances'),
    ('Desk Lamp', 39.99, 'Furniture'),
    ('Keyboard Mechanical', 89.99, 'Electronics'),
    ('Monitor 27"', 399.99, 'Electronics'),
    ('Webcam HD', 59.99, 'Electronics');

-- Copy data to data warehouse table (simulating ETL)
INSERT INTO products_dw (id, name, price, category, created_at)
SELECT id, name, price, category, created_at FROM products;

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    order_date DATE NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL
);

-- Create orders table for data warehouse
CREATE TABLE IF NOT EXISTS orders_dw (
    id INTEGER PRIMARY KEY,
    product_id INTEGER,
    quantity INTEGER NOT NULL,
    order_date DATE NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL
);

-- Insert sample orders
INSERT INTO orders (product_id, quantity, order_date, total_amount) VALUES
    (1, 2, '2024-01-15', 2599.98),
    (2, 5, '2024-01-15', 149.95),
    (3, 1, '2024-01-16', 249.99),
    (4, 1, '2024-01-16', 599.99),
    (5, 10, '2024-01-17', 199.90),
    (6, 2, '2024-01-17', 159.98),
    (7, 3, '2024-01-18', 119.97),
    (8, 1, '2024-01-18', 89.99),
    (9, 2, '2024-01-19', 799.98),
    (10, 4, '2024-01-19', 239.96);

-- Copy most orders to DW (but skip one to create a discrepancy)
INSERT INTO orders_dw (id, product_id, quantity, order_date, total_amount)
SELECT id, product_id, quantity, order_date, total_amount 
FROM orders 
WHERE id != 10; -- Intentionally skip order 10 to create a mismatch

-- Create customer table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    city VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create customer table for data warehouse
CREATE TABLE IF NOT EXISTS customers_dw (
    id INTEGER PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    city VARCHAR(100),
    created_at TIMESTAMP
);

-- Insert sample customers
INSERT INTO customers (name, email, city) VALUES
    ('John Doe', 'john.doe@example.com', 'New York'),
    ('Jane Smith', 'jane.smith@example.com', 'Los Angeles'),
    ('Bob Johnson', 'bob.johnson@example.com', 'Chicago'),
    ('Alice Brown', 'alice.brown@example.com', 'Houston'),
    ('Charlie Wilson', 'charlie.wilson@example.com', 'Phoenix');

-- Copy all customers to DW
INSERT INTO customers_dw (id, name, email, city, created_at)
SELECT id, name, email, city, created_at FROM customers;