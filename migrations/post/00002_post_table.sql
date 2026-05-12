-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

CREATE TABLE IF NOT EXISTS post.article (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_uuid UUID NOT NULL,
    username TEXT NOT NULL,
    description TEXT NOT NULL,
    media TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );

CREATE INDEX idx_article_user_uuid ON post.article(user_uuid);
CREATE INDEX idx_article_uuid_user_uuid ON post.article(uuid, user_uuid);

-- +goose Down
DROP INDEX IF EXISTS idx_article_user_uuid;
DROP INDEX IF EXISTS idx_article_uuid_user_uuid;

DROP TABLE IF EXISTS post.article;

DROP EXTENSION IF EXISTS pg_uuidv7;
