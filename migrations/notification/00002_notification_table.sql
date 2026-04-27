-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

CREATE TABLE IF NOT EXISTS notification.processed_event (
    uuid UUID PRIMARY KEY,
    event JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'processing',
    updated_at TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE notification.processed_event ADD CONSTRAINT chk_processed_event_status CHECK ( status IN ('processing', 'success') );

CREATE INDEX idx_processed_event_polling ON notification.processed_event (status, updated_at) WHERE status IN ('processing');

-- +goose Down
DROP INDEX IF EXISTS idx_processed_event_polling;

ALTER TABLE notification.processed_event DROP CONSTRAINT IF EXISTS chk_processed_event_status;

DROP TABLE IF EXISTS notification.processed_event;

DROP EXTENSION IF EXISTS pg_uuidv7;
