CREATE TABLE permissions
(
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(255) NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    category    VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT NOW(),
    deleted_by  VARCHAR(255)
)