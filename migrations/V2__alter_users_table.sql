ALTER TABLE users
    ADD COLUMN first_name TEXT NOT NULL,
    ADD COLUMN last_name TEXT NOT NULL,
    ADD COLUMN address_line1 TEXT NOT NULL,
    ADD COLUMN address_line2 TEXT,
    ADD COLUMN city TEXT NOT NULL,
    ADD COLUMN postal_code TEXT NOT NULL,
    ADD COLUMN country TEXT NOT NULL,
    ADD COLUMN phone_number TEXT NOT NULL,
    ADD COLUMN payment_method_id TEXT,
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
