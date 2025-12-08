CREATE TABLE breakdown_resources
(
    id           SERIAL PRIMARY KEY NOT NULL,
    breakdown_id BIGINT             NOT NULL REFERENCES breakdowns (id),
    image_path   VARCHAR(255),
    video_path   VARCHAR(255)
)