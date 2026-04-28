-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

CREATE TABLE IF NOT EXISTS post.article (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    username TEXT NOT NULL,
    description TEXT NOT NULL,
    media []BYTE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);

-- +goose Down
DROP INDEX IF EXISTS auth.idx_users_email;

DROP TABLE IF EXISTS post.article;

DROP EXTENSION IF EXISTS pg_uuidv7;
