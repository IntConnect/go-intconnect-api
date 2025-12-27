CREATE TABLE check_sheets
(
    id                               SERIAL PRIMARY KEY,
    check_sheet_document_template_id BIGINT REFERENCES check_sheet_document_templates (id) NOT NULL,
    reported_by                      BIGINT                                                NOT NULL REFERENCES users (id),
    verified_by                      BIGINT REFERENCES users (id),
    timestamp                        TIMESTAMP                                             NOT NULL DEFAULT current_timestamp,
    note                             TEXT,
    status                           VARCHAR(255)                                          NOT NULL DEFAULT 'Draft',
    created_at                       TIMESTAMP                                             NOT NULL DEFAULT current_timestamp,
    created_by                       VARCHAR(255),
    updated_at                       TIMESTAMP                                             NOT NULL DEFAULT current_timestamp,
    updated_by                       VARCHAR(255),
    deleted_at                       TIMESTAMP,
    deleted_by                       VARCHAR(255)
)