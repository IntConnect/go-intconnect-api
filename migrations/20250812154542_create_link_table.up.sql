CREATE TABLE links
(
    id             SERIAL PRIMARY KEY,
    pipeline_id    BIGINT NOT NULL REFERENCES pipelines (id),
    source_node_id BIGINT NOT NULL REFERENCES pipeline_nodes (id),
    target_node_id BIGINT NOT NULL REFERENCES pipeline_nodes (id)
);