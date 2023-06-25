BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS metrics (
        id SERIAL PRIMARY KEY,
        metric_name VARCHAR(100) UNIQUE,
        metric_type VARCHAR(7),
        metric_value DOUBLE PRECISION,
        metric_delta INT);

CREATE INDEX IF NOT EXISTS metric_name_idx ON metrics USING hash (metric_name);

COMMIT;