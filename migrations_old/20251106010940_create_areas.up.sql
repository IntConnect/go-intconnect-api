CREATE TABLE IF NOT EXISTS areas
(
    id         UUID         NOT NULL DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    position_x DECIMAL(10, 3),
    position_y DECIMAL(10, 3),
    position_z DECIMAL(10, 3),

    created_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    deleted_by VARCHAR(255),
    
    PRIMARY KEY (id)
);

