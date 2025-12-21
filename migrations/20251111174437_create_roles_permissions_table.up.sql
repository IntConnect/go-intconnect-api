CREATE TABLE roles_permissions
(
    id            SERIAL PRIMARY KEY,
    role_id       BIGINT NOT NULL REFERENCES roles (id),
    permission_id BIGINT NOT NULL REFERENCES permissions (id)
)