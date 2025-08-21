CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS income_item (
    id uuid primary key DEFAULT uuid_generate_v4(),
    name text NOT NULL
)