CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS subscriptions (
    id uuid primary key DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    type sub_status NOT NULL,
    duration_months INT NULL CHECK (duration_months IS NULL OR duration_months > 0),
    sessions_count INT NULL CHECK (sessions_count IS NULL OR sessions_count > 0),
    validity_months INT NULL CHECK (validity_months IS NULL OR validity_months > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);