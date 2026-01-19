CREATE TABLE modbus_servers
(
    id         SERIAL PRIMARY KEY,
    host_name  VARCHAR(255) NOT NULL,
    port       VARCHAR(255) NOT NULL,
    slave_id   VARCHAR(255) NOT NULL,
    timeout    INT          NOT NULL,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
)