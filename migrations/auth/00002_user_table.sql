-- +goose Up
CREATE TABLE IF NOT EXISTS auth.users (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);

-- +goose Down
DROP INDEX IF EXISTS auth.idx_users_email;

DROP TABLE IF EXISTS auth.users;
