CREATE TABLE IF NOT EXISTS semesters (
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subject      VARCHAR(255) NOT NULL,
    actual_grade SMALLINT NOT NULL CHECK (actual_grade BETWEEN 2 AND 5),
    parameters   JSONB NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_semesters_user_id ON semesters(user_id);
