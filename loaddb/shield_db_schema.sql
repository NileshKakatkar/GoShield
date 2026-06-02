CREATE TABLE serverinfo (
    uniq_id SERIAL UNIQUE,
    host_address TEXT NOT NULL,
    port INT,
    method TEXT,
    status TEXT
);