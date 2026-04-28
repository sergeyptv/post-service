-- +goose NO TRANSACTION

-- +goose Up
CREATE SCHEMA IF NOT EXISTS post;

-- +goose Down
DROP SCHEMA IF EXISTS post;
