CREATE TABLE check_sheet_values
(
    id             SERIAL PRIMARY KEY,
    check_sheet_id BIGINT       NOT NULL REFERENCES check_sheets (id),
    parameter_id   BIGINT       NOT NULL REFERENCES parameters (id),
    timestamp      VARCHAR(255) NOT NULL,
    value          TEXT         NOT NULL
)