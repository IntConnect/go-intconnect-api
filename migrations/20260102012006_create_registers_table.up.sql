CREATE TABLE registers
(
    id               SERIAL PRIMARY KEY,
    machine_id       BIGINT       NOT NULL REFERENCES machines (id),
    modbus_server_id BIGINT       NOT NULL REFERENCES modbus_servers (id),
    memory_location  VARCHAR(255) NOT NULL,
    name             VARCHAR(255) NOT NULL,
    description      VARCHAR(255) NOT NULL,
    data_type        VARCHAR(255) NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    created_by       VARCHAR(255),
    updated_at       TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP,
    deleted_by       VARCHAR(255)
)
