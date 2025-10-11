CREATE TABLE pipelines
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    created_by  VARCHAR(100),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    is_active   BOOLEAN   DEFAULT TRUE,
    -- optional: simpan snapshot seluruh pipeline sebagai JSON
    config      JSONB
);