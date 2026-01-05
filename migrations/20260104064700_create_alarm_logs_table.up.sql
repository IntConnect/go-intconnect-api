CREATE TABLE alarm_logs
(
    id              SERIAL PRIMARY KEY,
    parameter_id    BIGINT REFERENCES parameters (id) NOT NULL,
    acknowledged_by BIGINT REFERENCES users (id),
    value           FLOAT                             NOT NULL,
    type            VARCHAR(255)                      NOT NULL,
    is_active       BOOLEAN                           NOT NULL DEFAULT TRUE,
    status          VARCHAR(255)                      NOT NULL DEFAULT 'Open',
    note            VARCHAR(255),
    acknowledged_at TIMESTAMP,
    resolved_at     TIMESTAMP,
    created_at      TIMESTAMP                                  DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP                                  DEFAULT CURRENT_TIMESTAMP
)