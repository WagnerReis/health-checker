CREATE TABLE IF NOT EXISTS health_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    monitor_id UUID NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status_code INTEGER,
    response_time_ms INTEGER NOT NULL,
    is_success BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_health_checks_monitor_checked
    ON health_checks (monitor_id, checked_at DESC);