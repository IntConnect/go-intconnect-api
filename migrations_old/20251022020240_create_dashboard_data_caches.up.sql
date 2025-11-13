CREATE TABLE dashboard_data_cache
(
    id         SERIAL PRIMARY KEY,
    source_key VARCHAR(100)             NOT NULL,
    cache_key  VARCHAR(255) UNIQUE      NOT NULL, -- Hash dari source_key + params
    data       JSONB                    NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);