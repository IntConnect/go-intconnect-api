CREATE TABLE parameter_operations
(
    id           SERIAL PRIMARY KEY NOT NULL,
    type         VARCHAR(255) DEFAULT '',
    value        FLOAT,
    sequence     INT,
    parameter_id BIGINT             NOT NULL REFERENCES parameters (id)
)