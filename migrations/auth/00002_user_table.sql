-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

CREATE TABLE IF NOT EXISTS auth.users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    username TEXT NOT NULL UNIQUE,
    passHash TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);

-- +goose Down
DROP INDEX IF EXISTS auth.idx_users_email;

DROP TABLE IF EXISTS auth.users;

DROP EXTENSION IF EXISTS pg_uuidv7;
