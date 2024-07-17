CREATE SCHEMA IF NOT EXISTS messaggio;

CREATE TABLE messaggio.messages (
     id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP
);