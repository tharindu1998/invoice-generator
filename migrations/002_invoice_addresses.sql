ALTER TABLE invoices
    ADD COLUMN customer_address TEXT NOT NULL DEFAULT '',
    ADD COLUMN seller_name      VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN seller_address   TEXT NOT NULL DEFAULT '';
