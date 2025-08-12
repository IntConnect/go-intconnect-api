CREATE TABLE pipeline_nodes
(
    id          SERIAL PRIMARY KEY,
    pipeline_id BIGINT NOT NULL REFERENCES pipelines (id),
    node_id     BIGINT NOT NULL REFERENCES nodes (id),
    type        TEXT   NOT NULL,               -- contoh: "mqtt-in", "rs232-in", "http-out", "function"
    name        TEXT   NOT NULL,
    description TEXT,
    config      JSONB     DEFAULT '{}'::jsonb, -- konfigurasi node
    position_x  INT       DEFAULT 0,
    position_y  INT       DEFAULT 0,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
);
