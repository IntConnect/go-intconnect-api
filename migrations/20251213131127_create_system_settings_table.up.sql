CREATE TABLE system_settings
(
    id          SERIAL PRIMARY KEY,
    key         VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    value       JSONB        NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    created_by  VARCHAR(255),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP,
    deleted_by  VARCHAR(255)
);

