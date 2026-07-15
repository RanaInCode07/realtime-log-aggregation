CREATE TABLE IF NOT EXISTS system_logs(
    id UUID PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    level varchar(10) NOT NULL,
    service varchar(50) NOT NULL,
    message TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON system_logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_logs_service_level ON system_logs(service, level);