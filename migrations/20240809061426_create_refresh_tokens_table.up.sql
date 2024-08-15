CREATE TABLE refresh_tokens(
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    token_id UUID NOT NULL,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
