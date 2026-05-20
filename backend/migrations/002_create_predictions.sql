CREATE TABLE IF NOT EXISTS predictions (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parameters      JSONB NOT NULL,
    predicted_grade SMALLINT NOT NULL CHECK (predicted_grade BETWEEN 2 AND 5),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_predictions_user_id ON predictions(user_id);
