CREATE TABLE machine_documents
(
    id          SERIAL       NOT NULL,
    machine_id  BIGINT       NOT NULL REFERENCES machines (id),
    code        VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    file_path   VARCHAR(255) NOT NULL,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by  VARCHAR(255)
)