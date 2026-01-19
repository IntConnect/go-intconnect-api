CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE telemetries
(
    id           BIGSERIAL NOT NULL,
    parameter_id BIGINT    NOT NULL REFERENCES parameters (id),
    value        FLOAT,
    timestamp    TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

SELECT create_hypertable(
               'telemetries',
               'timestamp',
               'id',
               number_partitions => 1
       )