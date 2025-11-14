CREATE TABLE parameters
(
    id          SERIAL PRIMARY KEY,
    machine_id  BIGINT       NOT NULL REFERENCES machines (id),
    name        VARCHAR(255) NOT NULL,
    code        VARCHAR(255) NOT NULL,
    unit        VARCHAR(255) NOT NULL,
    min_value   FLOAT,
    max_value   FLOAT,
    description TEXT,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
)