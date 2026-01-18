CREATE TABLE dashboard_widgets
(
    id         SERIAL PRIMARY KEY,
    machine_id BIGINT       NOT NULL REFERENCES machines (id),
    code       VARCHAR(255) NOT NULL,
    layout     JSONB        NOT NULL DEFAULT '{}'::jsonb,
    config     JSONB        NOT NULL DEFAULT '{}'::jsonb
);