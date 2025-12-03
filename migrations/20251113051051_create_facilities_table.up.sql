CREATE TYPE facility_status_enum AS ENUM ('Active', 'Maintenance', 'Archived');


CREATE TABLE facilities
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255)         NOT NULL,
    code           VARCHAR(255)         NOT NULL UNIQUE,
    description    TEXT,
    location       VARCHAR(255),
    status         facility_status_enum NOT NULL,
    thumbnail_path VARCHAR(255)         NOT NULL,
    metadata       JSONB     DEFAULT '{}'::jsonb,
    created_by     VARCHAR(255),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by     VARCHAR(255)
)