CREATE TABLE IF NOT EXISTS events (
    event_id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    action VARCHAR(50) NOT NULL,
    action_object_id BIGINT,
    action_object_type VARCHAR(50),
    timestamp TIMESTAMPTZ NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_events_user_timestamp ON events (user_id, timestamp);
