BEGIN;

CREATE TABLE IF NOT EXISTS titles (
    title VARCHAR(10) NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS coins (
    title         VARCHAR(10) NOT NULL REFERENCES titles(title),
    rate          REAL NOT NULL CHECK (rate >= 0),
    creation_time TIMESTAMPTZ NOT NULL
    );

INSERT INTO titles (title) VALUES
    ('BTC')
ON CONFLICT (title) DO NOTHING;

END;