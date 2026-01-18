CREATE TABLE users
(
    id          SERIAL PRIMARY KEY,
    role_id     BIGINT REFERENCES roles (id) NOT NULL,
    username    VARCHAR(255)                 NOT NULL UNIQUE,
    name        VARCHAR(255)                 NOT NULL UNIQUE,
    email       VARCHAR(255)                 NOT NULL UNIQUE,
    password    VARCHAR(255)                 NOT NULL,
    avatar_path VARCHAR(255)                 NOT NULL,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
)