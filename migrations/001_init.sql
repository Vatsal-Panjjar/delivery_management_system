CREATE EXTENSION IF NOT EXISTS pgcrypto;


CREATE TABLE IF NOT EXISTS users (
id SERIAL PRIMARY KEY,
name TEXT NOT NULL,
email TEXT UNIQUE NOT NULL,
password_hash TEXT NOT NULL,
role TEXT NOT NULL DEFAULT 'customer',
created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);


CREATE TABLE IF NOT EXISTS deliveries (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
customer_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
courier_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
pickup_address TEXT NOT NULL,
dropoff_address TEXT NOT NULL,
status TEXT NOT NULL DEFAULT 'pending',
price_cents INTEGER NOT NULL DEFAULT 0,
created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);


CREATE TABLE IF NOT EXISTS delivery_history (
id SERIAL PRIMARY KEY,
delivery_id UUID REFERENCES deliveries(id) ON DELETE CASCADE,
status TEXT NOT NULL,
changed_by INTEGER,
note TEXT,
changed_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);


CREATE INDEX IF NOT EXISTS idx_deliveries_status ON deliveries(status);
