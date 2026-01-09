CREATE TABLE machine_documents
(
    id          SERIAL PRIMARY KEY,
    machine_id  BIGINT       NOT NULL REFERENCES machines (id),
    code        VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    file_path   VARCHAR(255) NOT NULL
)