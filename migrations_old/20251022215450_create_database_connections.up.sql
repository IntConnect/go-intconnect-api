CREATE TABLE database_connections
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    database_type VARCHAR(100) NOT NULL,
    database_name VARCHAR(100) NOT NULL,
    description   TEXT,
    config        JSONB,
    created_by    VARCHAR(255),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by    VARCHAR(255),
    deleted_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by    VARCHAR(255)
)