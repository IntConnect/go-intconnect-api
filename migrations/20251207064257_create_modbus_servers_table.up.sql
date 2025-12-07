CREATE TABLE modbus_servers
(
    id         SERIAL PRIMARY KEY NOT NULL,
    ip_address VARCHAR(255)       NOT NULL,
    port       VARCHAR(255)       NOT NULL,
    slave_id   VARCHAR(255)       NOT NULL,
    timeout    INT                NOT NULL,
    is_active  BOOLEAN                     DEFAULT FALSE,
    created_at TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
)