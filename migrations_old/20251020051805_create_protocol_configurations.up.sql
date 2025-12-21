CREATE TABLE protocol_configurations
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    protocol    VARCHAR(50)  NOT NULL,
    description TEXT,
    config      JSONB        NOT NULL,
    is_active   BOOLEAN   DEFAULT TRUE,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
);