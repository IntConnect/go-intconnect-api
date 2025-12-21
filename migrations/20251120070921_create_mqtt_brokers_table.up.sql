CREATE TABLE mqtt_brokers
(
    id         SERIAL PRIMARY KEY,
    host_name  VARCHAR(100) NOT NULL,
    mqtt_port  VARCHAR(100) NOT NULL,
    ws_port    VARCHAR(100) NOT NULL,
    username   VARCHAR(255),
    password   VARCHAR(255),
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP   NOT NULL DEFAULT current_timestamp,
    created_by VARCHAR(255),
    updated_at TIMESTAMP   NOT NULL DEFAULT current_timestamp,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP,
    deleted_by VARCHAR(255)
);

