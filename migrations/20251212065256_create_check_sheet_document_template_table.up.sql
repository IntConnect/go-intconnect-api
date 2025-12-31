CREATE TABLE check_sheet_document_templates
(
    id              SERIAL PRIMARY KEY,
    machine_id      BIGINT REFERENCES machines (id) NOT NULL,
    name            VARCHAR(255)                    NOT NULL,
    no              VARCHAR(255)                    NOT NULL,
    description     VARCHAR(255),
    category        VARCHAR(255)                    NOT NULL,
    interval        INT                             NOT NULL DEFAULT 1,
    interval_type   VARCHAR(255)                    NOT NULL DEFAULT 'Hours',
    rotation_type   VARCHAR                         NOT NULL DEFAULT 'Daily',
    revision_number INT                             NOT NULL DEFAULT 0,
    effective_date  DATE                            NOT NULL,
    created_by      VARCHAR(255),
    created_at      TIMESTAMP                                DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP                                DEFAULT CURRENT_TIMESTAMP,
    updated_by      VARCHAR(255),
    deleted_at      TIMESTAMP                                DEFAULT CURRENT_TIMESTAMP,
    deleted_by      VARCHAR(255)
)