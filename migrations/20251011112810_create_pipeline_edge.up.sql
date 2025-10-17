CREATE TABLE pipeline_edges
(
    id             SERIAL PRIMARY KEY,
    pipeline_id    INT REFERENCES pipelines (id) ON DELETE CASCADE,
    source_node_id BIGINT NOT NULL REFERENCES pipeline_nodes (id),
    target_node_id BIGINT NOT NULL REFERENCES pipeline_nodes (id),
    data           JSONB, -- optional: misalnya tipe koneksi atau kondisi
    created_at     TIMESTAMP DEFAULT NOW()
);