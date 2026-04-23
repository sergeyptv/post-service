-- +goose NO TRANSACTION

-- +goose Up
CREATE SCHEMA IF NOT EXISTS notification;

-- +goose Down
DROP SCHEMA IF EXISTS notification;
