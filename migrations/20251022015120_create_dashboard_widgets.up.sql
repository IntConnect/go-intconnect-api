CREATE TABLE dashboard_widgets
(
    id               SERIAL PRIMARY KEY,
    dashboard_id     BIGINT       NOT NULL REFERENCES dashboard_configs (id),
    widget_key       VARCHAR(100) NOT NULL, -- Unique identifier untuk grid layout (i)
    widget_type      VARCHAR(100) NOT NULL,
    title            VARCHAR(255) NOT NULL,
    data_source_key  VARCHAR(100) NOT NULL, -- Reference ke data source

    -- Grid Layout Position & Size
    position_x       INTEGER      NOT NULL    DEFAULT 0,
    position_y       INTEGER      NOT NULL    DEFAULT 0,
    width            INTEGER      NOT NULL    DEFAULT 6,
    height           INTEGER      NOT NULL    DEFAULT 6,

    -- Widget Specific Configuration
    chart_config     JSONB,                 -- ApexCharts options (colors, legends, etc)
    filters          JSONB,                 -- Filter/parameters untuk data
    refresh_interval INTEGER,               -- Auto refresh dalam detik (null = no refresh)

    -- Display Options
    is_visible       BOOLEAN                  DEFAULT TRUE,
    display_order    INTEGER                  DEFAULT 0,

    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_widget_key UNIQUE (dashboard_id, widget_key)
);