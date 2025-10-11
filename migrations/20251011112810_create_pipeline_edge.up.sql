CREATE TABLE pipeline_edges
(
    id             SERIAL PRIMARY KEY,
    pipeline_id    INT REFERENCES pipelines (id) ON DELETE CASCADE,
    edge_id        VARCHAR(50) NOT NULL,
    source_node_id VARCHAR(50) NOT NULL,
    target_node_id VARCHAR(50) NOT NULL,
    data           JSONB, -- optional: misalnya tipe koneksi atau kondisi
    created_at     TIMESTAMP DEFAULT NOW()
);