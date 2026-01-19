CREATE TABLE check_sheet_document_templates
(
    id              SERIAL PRIMARY KEY,
    machine_id      BIGINT REFERENCES machines (id) NOT NULL,
    name            VARCHAR(255)                    NOT NULL,
    no              VARCHAR(255)                    NOT NULL,
    description     VARCHAR(255),
    category        VARCHAR(255)                    NOT NULL,
    interval        INT                             NOT NULL DEFAULT 1,
    interval_type   VARCHAR(255)                    NOT NULL DEFAULT 'Hour',
    rotation_type   VARCHAR                         NOT NULL DEFAULT 'Daily',
    revision_number INT                             NOT NULL DEFAULT 0,
    effective_date  DATE                            NOT NULL,
    starting_hour   TIME                            NOT NULL,
    created_by      VARCHAR(255),
    created_at      TIMESTAMP                                DEFAULT NOW(),
    updated_at      TIMESTAMP                                DEFAULT NOW(),
    updated_by      VARCHAR(255),
    deleted_at      TIMESTAMP                                DEFAULT NOW(),
    deleted_by      VARCHAR(255)
)