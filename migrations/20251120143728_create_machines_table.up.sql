CREATE TABLE machines
(
    id             SERIAL PRIMARY KEY,
    facility_id    BIGINT       NOT NULL REFERENCES facilities (id),
    name           VARCHAR(255) NOT NULL,
    code           VARCHAR(255) NOT NULL,
    description    VARCHAR(255) NOT NULL,
    camera_x       FLOAT        NOT NULL DEFAULT 0,
    camera_y       FLOAT        NOT NULL DEFAULT 0,
    camera_z       FLOAT        NOT NULL DEFAULT 0,
    thumbnail_path VARCHAR(255) NOT NULL,
    model_path     VARCHAR(255) NOT NULL,
    created_by     VARCHAR(255),
    created_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    deleted_by     VARCHAR(255)
)