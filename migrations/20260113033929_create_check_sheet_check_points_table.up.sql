CREATE TABLE check_sheet_check_points
(
    id             SERIAL PRIMARY KEY,
    check_sheet_id BIGINT       NOT NULL REFERENCES check_sheets (id),
    parameter_id   BIGINT       NOT NULL REFERENCES parameters (id),
    name           VARCHAR(255) NOT NULL DEFAULT '',
    created_by     VARCHAR(255),
    created_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255)
)