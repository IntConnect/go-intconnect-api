CREATE TABLE telemetries
(
    id           SERIAL PRIMARY KEY,
    parameter_id BIGINT NOT NULL REFERENCES parameters (id),
    value        FLOAT,
    timestamp    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)