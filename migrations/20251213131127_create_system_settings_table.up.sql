CREATE TABLE system_settings
(
    id          SERIAL PRIMARY KEY,
    key         VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    value       JSONB        NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP    NOT NULL DEFAULT current_timestamp,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP,
    deleted_by  VARCHAR(255)
);

