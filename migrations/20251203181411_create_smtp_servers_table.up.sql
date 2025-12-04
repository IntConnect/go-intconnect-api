CREATE TABLE smtp_servers
(
    id           SERIAL PRIMARY KEY NOT NULL,
    host         VARCHAR(255)       NOT NULL,
    port         VARCHAR(255)       NOT NULL,
    username     VARCHAR(255)       NOT NULL,
    password     VARCHAR(255)       NOT NULL,
    mail_address VARCHAR(255)       NOT NULL,
    mail_name    VARCHAR(255)       NOT NULL,
    is_active    BOOLEAN                     DEFAULT TRUE NOT NULL,
    created_at   TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    created_by   VARCHAR(255),
    updated_at   TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    updated_by   VARCHAR(255),
    deleted_at   TIMESTAMPTZ,
    deleted_by   VARCHAR(255)
)