CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS students (
    id uuid primary key DEFAULT uuid_generate_v4(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    full_name text NOT NULL,
    gender gender NOT NULL DEFAULT 'мужчина',
    phoneNumber text,
    parentNumber text,
    status student_status NOT NULL DEFAULT 'активный',
    note text,
    version integer NOT NULL DEFAULT 1
)