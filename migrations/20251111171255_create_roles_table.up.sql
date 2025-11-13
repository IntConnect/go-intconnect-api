CREATE TABLE roles
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    description    TEXT         NOT NULL,
    is_system_role BOOLEAN      NOT NULL DEFAULT FALSE,
    created_by     VARCHAR(255),
    created_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    deleted_by     VARCHAR(255)
)