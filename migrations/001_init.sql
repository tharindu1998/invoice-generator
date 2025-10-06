CREATE TABLE customers (
    mobile VARCHAR(20) PRIMARY KEY,
    name   VARCHAR(100) NOT NULL,
    email  VARCHAR(100)
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL
);

CREATE TABLE invoices (
    id SERIAL PRIMARY KEY,
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    customer_mobile VARCHAR(20),
    customer_name   VARCHAR(100) NOT NULL,
    customer_email  VARCHAR(100),
    date DATE NOT NULL,
    due_date DATE,
    total_amount DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (customer_mobile) REFERENCES customers(mobile)
);

CREATE TABLE invoice_items (
    id SERIAL PRIMARY KEY,
    invoice_id INT REFERENCES invoices(id),
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total DECIMAL(10,2) NOT NULL
);
