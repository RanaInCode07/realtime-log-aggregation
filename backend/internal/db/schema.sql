CREATE TABLE IF NOT EXISTS system_logs(
    id UUID NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    level varchar(10) NOT NULL,
    service varchar(50) NOT NULL,
    message TEXT NOT NULL,
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

CREATE TABLE IF NOT EXISTS system_logs_default 
PARTITION OF system_logs DEFAULT;

CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON system_logs(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_logs_service_level ON system_logs(service, level);