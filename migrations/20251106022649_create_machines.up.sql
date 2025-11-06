CREATE TABLE IF NOT EXISTS machines
(
    id         SERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(150)       NOT NULL,
    parameter  VARCHAR(50),
    photo      VARCHAR(255),
    asset_3d   VARCHAR(255),
    camera_x   DECIMAL(10, 3) DEFAULT 0,
    camera_y   DECIMAL(10, 3) DEFAULT 0,
    camera_z   DECIMAL(10, 3) DEFAULT 0,
    camera_fov DECIMAL(10, 3) DEFAULT 50,
    area_id    UUID               NOT NULL,

    created_at TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    deleted_by VARCHAR(255),

    PRIMARY KEY (id)
);
