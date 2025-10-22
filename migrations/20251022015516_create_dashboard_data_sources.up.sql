CREATE TABLE dashboard_data_sources
(
    id            SERIAL PRIMARY KEY,
    source_key    VARCHAR(100) UNIQUE NOT NULL,
    source_name   VARCHAR(255)        NOT NULL,
    description   TEXT,

    -- Static data (untuk development/testing)
    static_data   JSONB,

    -- API Configuration
    api_endpoint  VARCHAR(500),
    api_method    VARCHAR(10)              DEFAULT 'GET',
    api_headers   JSONB,

    -- Database Query
    sql_query     TEXT,

    -- Cache Configuration
    cache_enabled BOOLEAN                  DEFAULT FALSE,
    cache_ttl     INTEGER                  DEFAULT 300, -- seconds

    -- Metadata
    category      VARCHAR(100),
    is_active     BOOLEAN                  DEFAULT TRUE,

    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);