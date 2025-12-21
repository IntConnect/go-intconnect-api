CREATE TABLE parameter_operations
(
    id           SERIAL PRIMARY KEY,
    parameter_id BIGINT       NOT NULL REFERENCES parameters (id),
    type         VARCHAR(255) NOT NULL,
    value        FLOAT        NOT NULL,
    sequence     INT          NOT NULL
)