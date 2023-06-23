CREATE USER evgeny
    PASSWORD 'evgeny';

CREATE DATABASE evgeny_trefilov
    OWNER 'evgeny'
    ENCODING 'UTF-8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';

CREATE TABLE IF NOT EXISTS metrics (
        id SERIAL PRIMARY KEY,
        metric_name VARCHAR(100) UNIQUE,
        metric_type VARCHAR(7),
        metric_value DOUBLE PRECISION,
        metric_delta INT);