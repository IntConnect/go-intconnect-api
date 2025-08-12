CREATE TABLE nodes
(
    id             SERIAL PRIMARY KEY NOT NULL,
    type           VARCHAR(255)       NOT NULL,
    label          VARCHAR(255)       NOT NULL,
    description    VARCHAR(255)       NOT NULL,
    help_text      TEXT               NOT NULL,
    color          VARCHAR(255)       NOT NULL,
    icon           VARCHAR(255)       NOT NULL,
    component_name VARCHAR(255)       NOT NULL,
    default_config JSONB DEFAULT '{}'::jsonb
)