CREATE INDEX idx_teachers_full_name_gin ON teachers USING GIN (to_tsvector('russian', full_name));
CREATE INDEX idx_teachers_status ON teachers(status);
