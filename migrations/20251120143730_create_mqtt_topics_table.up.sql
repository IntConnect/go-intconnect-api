CREATE TABLE mqtt_topics
(
    id             SERIAL PRIMARY KEY NOT NULL,
    machine_id     BIGINT             NOT NULL REFERENCES machines (id) UNIQUE,
    mqtt_broker_id BIGINT             NOT NULL REFERENCES mqtt_brokers (id),
    name           VARCHAR(255)       NOT NULL,
    qos            INT                NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    created_by     VARCHAR(255),
    updated_at     TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMPTZ,
    deleted_by     VARCHAR(255)
)