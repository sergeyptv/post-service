-- +goose Up
CREATE SCHEMA IF NOT EXISTS outbox;

CREATE TABLE IF NOT EXISTS outbox.event (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    version TEXT NOT NULL,
    user_uuid UUID NOT NULL ,
    username TEXT NOT NULL,
    user_email TEXT NOT NULL,
    registered_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
);

ALTER TABLE outbox.event ADD CONSTRAINT chk_event_status CHECK ( status IN ('pending', 'processing', 'sent', 'failed') );

CREATE INDEX idx_event_status ON outbox.event (status) WHERE status = 'pending';


-- +goose Down
DROP INDEX IF EXISTS idx_event_status;

ALTER TABLE outbox.event DROP CONSTRAINT IF EXISTS chk_event_status;

DROP TABLE IF EXISTS outbox.event;

DROP SCHEMA IF EXISTS outbox;
