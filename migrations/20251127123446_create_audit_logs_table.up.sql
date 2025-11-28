CREATE TABLE audit_logs
(
    id          SERIAL       NOT NULL PRIMARY KEY,
    user_id     BIGINT REFERENCES users (id),
    action      VARCHAR(255) NOT NULL,
    feature     VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    before      JSONB     DEFAULT '{}'::jsonb,
    after       JSONB     DEFAULT '{}'::jsonb,
    ip_address  VARCHAR(255) NOT NULL,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
)