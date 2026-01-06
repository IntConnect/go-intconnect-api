CREATE TABLE processed_parameter_sequences
(
    id                  SERIAL PRIMARY KEY,
    parent_parameter_id BIGINT       NOT NULL REFERENCES parameters (id),
    parameter_id        BIGINT       NOT NULL REFERENCES parameters (id),
    sequence            INT          NOT NULL DEFAULT 1,
    type                VARCHAR(255) NOT NULL

)