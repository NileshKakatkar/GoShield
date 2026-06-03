CREATE TABLE serverinfo (
    uniq_id SERIAL UNIQUE,
    host_address TEXT NOT NULL,
    port INT,
    method TEXT,
    status TEXT
);

CREATE TABLE admin(
    uniq_id SERIAL UNIQUE,
    admin_name TEXT,
    admin_id TEXT,
    password TEXT
);

INSERT INTO admin (admin_name, admin_id, password) VALUES ('dev', 'admin', '$2a$10$vDV4RueUHDx43rrySODBhOznLua3wPdy5q.bpreeLzvpEa.oswczm');