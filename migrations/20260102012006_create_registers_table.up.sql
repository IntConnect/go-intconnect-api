CREATE TABLE registers
(
    id               SERIAL PRIMARY KEY,
    machine_id       BIGINT       NOT NULL REFERENCES machines (id),
    modbus_server_id BIGINT       NOT NULL REFERENCES modbus_servers (id),
    memory_location  VARCHAR(255) NOT NULL,
    name             VARCHAR(255) NOT NULL,
    description      VARCHAR(255) NOT NULL,
    data_type        VARCHAR(255) NOT NULL,
    unit             VARCHAR(255) NOT NULL,
    position_x       FLOAT        NOT NULL DEFAULT 0,
    position_y       FLOAT        NOT NULL DEFAULT 0,
    position_z       FLOAT        NOT NULL DEFAULT 0,
    rotation_x       FLOAT        NOT NULL DEFAULT 0,
    rotation_y       FLOAT        NOT NULL DEFAULT 0,
    rotation_z       FLOAT        NOT NULL DEFAULT 0,
    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    created_by       VARCHAR(255),
    updated_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP,
    deleted_by       VARCHAR(255)
)
