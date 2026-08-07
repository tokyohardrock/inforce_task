CREATE TABLE user_activity_stats (
    user_id BIGINT NOT NULL,
    time_bucket TIMESTAMPTZ NOT NULL,
    event_count INT NOT NULL,

    PRIMARY KEY (user_id, time_bucket)
);
