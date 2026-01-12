CREATE TABLE roles_permissions
(
    role_id       BIGINT NOT NULL REFERENCES roles (id),
    permission_id BIGINT NOT NULL REFERENCES permissions (id),
    UNIQUE (role_id, permission_id)
)