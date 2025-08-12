CREATE TABLE pipelines
(
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    is_active   BOOLEAN   DEFAULT true,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
);