CREATE TABLE nodes
(
    id             SERIAL PRIMARY KEY,
    type           VARCHAR(255) NOT NULL,
    name           VARCHAR(255) NOT NULL,
    label          VARCHAR(255) NOT NULL,
    description    VARCHAR(255) NOT NULL,
    help_text      TEXT         NOT NULL,
    color          VARCHAR(255) NOT NULL,
    icon           VARCHAR(255) NOT NULL,
    component_name VARCHAR(255) NOT NULL,
    default_config JSONB     DEFAULT '{}'::jsonb,
    created_by     VARCHAR(255),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by     VARCHAR(255)
)