CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS teachers (
    id uuid primary key DEFAULT uuid_generate_v4(),
    full_name text NOT NULL,
    gender gender NOT NULL DEFAULT 'мужчина',
    birth_date timestamp(0) with time zone,
    phone text,
    note text,
    status teacher_status NOT NULL DEFAULT 'активный',
    created_at timestamp(0) with time zone NOT NULL default now(),
    updated_at timestamp(0) with time zone NOT NULL default now()
);
