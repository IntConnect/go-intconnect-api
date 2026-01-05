

CREATE TABLE facilities
(
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255)         NOT NULL,
    code           VARCHAR(255)         NOT NULL UNIQUE,
    location       VARCHAR(255),
    description    TEXT,
    thumbnail_path VARCHAR(255)         NOT NULL,
    model_path     VARCHAR(255)         NOT NULL,
    position_x     FLOAT                NOT NULL DEFAULT 0,
    position_y     FLOAT                NOT NULL DEFAULT 0,
    position_z     FLOAT                NOT NULL DEFAULT 0,
    camera_x       FLOAT                NOT NULL DEFAULT 0,
    camera_y       FLOAT                NOT NULL DEFAULT 0,
    camera_z       FLOAT                NOT NULL DEFAULT 0,
    created_by     VARCHAR(255),
    created_at     TIMESTAMP                     DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP                     DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP,
    deleted_by     VARCHAR(255)
)