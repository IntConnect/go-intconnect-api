CREATE TABLE machine_documents
(
    id          SERIAL PRIMARY KEY,
    machine_id  BIGINT       NOT NULL REFERENCES machines (id),
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    file_path   VARCHAR(255) NOT NULL,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP,
    deleted_by  VARCHAR(255)
)