-- +goose Up
CREATE SCHEMA IF NOT EXISTS token;

CREATE TABLE IF NOT EXISTS token.storage (
    jti UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_uuid UUID NOT NULL UNIQUE,
    token TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS token.storage;

DROP SCHEMA IF EXISTS token;
