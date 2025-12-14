CREATE TABLE system_settings
(
    id          BIGSERIAL PRIMARY KEY,
    key         VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    value       JSONB        NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp,
    created_by  VARCHAR(255),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT current_timestamp,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMPTZ,
    deleted_by  VARCHAR(255)
);

