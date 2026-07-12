CREATE TABLE IF NOT EXISTS open_ports (
    id BIGSERIAL PRIMARY KEY,
    ip INET NOT NULL,
    port INTEGER NOT NULL,
    proto TEXT NOT NULL DEFAULT 'tcp',
    service TEXT NOT NULL DEFAULT '',
    banner TEXT NOT NULL DEFAULT '',
    first_seen TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (ip, port, proto)
);

CREATE INDEX IF NOT EXISTS idx_open_ports_last_seen ON open_ports (last_seen);
