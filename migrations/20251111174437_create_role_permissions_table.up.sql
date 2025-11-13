CREATE TABLE role_permissions
(
    id            SERIAL PRIMARY KEY,
    role_id       BIGINT  NOT NULL,
    permission_id BIGINT  NOT NULL,
    granted       BOOLEAN NOT NULL DEFAULT TRUE
)