CREATE TABLE report_document_templates
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    code             VARCHAR(255) NOT NULL,
    document_version INT          NOT NULL DEFAULT 0,
    created_by       VARCHAR(255),
    created_at       TIMESTAMP             DEFAULT NOW(),
    updated_at       TIMESTAMP             DEFAULT NOW(),
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP,
    deleted_by       VARCHAR(255)
)