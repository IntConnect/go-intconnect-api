CREATE TABLE machines
(
    id             SERIAL PRIMARY KEY,
    facility_id    BIGINT       NOT NULL REFERENCES facilities (id),
    mqtt_topic_id  BIGINT       NOT NULL REFERENCES mqtt_topics (id),
    name           VARCHAR(255) NOT NULL,
    code           VARCHAR(255) NOT NULL,
    description    VARCHAR(255) NOT NULL,
    model_path     VARCHAR(255) NOT NULL,
    model_offset_x FLOAT        NOT NULL,
    model_offset_y FLOAT        NOT NULL,
    model_offset_z FLOAT        NOT NULL,
    model_scale    FLOAT        NOT NULL,
    thumbnail_path VARCHAR(255) NOT NULL,
    created_by     VARCHAR(255),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by     VARCHAR(255)
)