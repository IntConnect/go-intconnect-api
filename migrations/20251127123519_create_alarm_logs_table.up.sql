CREATE TABLE alarm_logs
(
    id              SERIAL PRIMARY KEY,
    parameter_id    BIGINT REFERENCES parameters (id),
    value           FLOAT        NOT NULL,
    type            VARCHAR(255) NOT NULL,
    category        VARCHAR(255) NOT NULL,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    status          VARCHAR(255) NOT NULL,
    acknowledged_at TIMESTAMP,
    finished_at     TIMESTAMP,
    notes           VARCHAR(255),
    created_by      VARCHAR(255),
    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_by      VARCHAR(255),
    deleted_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    deleted_by      VARCHAR(255)
)