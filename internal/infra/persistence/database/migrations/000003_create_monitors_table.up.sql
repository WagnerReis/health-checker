CREATE TABLE IF NOT EXISTS monitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    method VARCHAR(255) NOT NULL DEFAULT 'GET',
    headers JSONB DEFAULT '{}',
    body TEXT,
    interval INT NOT NULL,
    expected_status_code INT,
    timeout INT NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'ACTIVE',
    deleted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_monitors_user_id ON monitors (user_id);
CREATE INDEX IF NOT EXISTS idx_monitors_name ON monitors (name);