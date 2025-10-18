CREATE TABLE pipeline_nodes
(
    id          SERIAL PRIMARY KEY,
    pipeline_id INT REFERENCES pipelines (id) ON DELETE CASCADE,
    node_id     BIGINT REFERENCES nodes (id),
    type        VARCHAR(50) NOT NULL, -- contoh: mqtt-in, json-parser, mqtt-out
    label       VARCHAR(100),
    position_x  FLOAT,
    position_y  FLOAT,
    config      JSONB,                -- semua konfigurasi unik node (topic, QoS, dll)
    description TEXT
);
