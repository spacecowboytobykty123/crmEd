CREATE INDEX IF NOT EXISTS teacher_name_idx ON teachers USING GIN (to_tsvector('simple', full_name));
