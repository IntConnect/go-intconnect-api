CREATE TABLE check_sheet_check_point_values
(
    id                        SERIAL PRIMARY KEY,
    check_sheet_check_point_id BIGINT       NOT NULL REFERENCES check_sheet_check_points (id),
    timestamp                 VARCHAR(255) NOT NULL,
    value                     TEXT         NOT NULL
)