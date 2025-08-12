CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by VARCHAR(255)
)