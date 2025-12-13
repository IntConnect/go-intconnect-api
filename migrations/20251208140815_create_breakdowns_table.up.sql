CREATE TABLE breakdowns
(
    id                     SERIAL PRIMARY KEY NOT NULL,
    machine_id             BIGINT             NOT NULL REFERENCES machines (id),
    reported_by            BIGINT             NOT NULL REFERENCES users (id),
    verified_by            BIGINT REFERENCES users (id),

    problem_identification TEXT,
    people_issue           TEXT               NOT NULL DEFAULT '',
    inspection_issue       TEXT               NOT NULL DEFAULT '',
    condition_issue        TEXT               NOT NULL DEFAULT '',
    action_taken           TEXT               NOT NULL DEFAULT '',
    parts_issue            TEXT               NOT NULL DEFAULT '',
    analysis_notes         TEXT               NOT NULL DEFAULT '',

    corrective_action      TEXT,
    preventive_action      TEXT,
    status                 VARCHAR(255)       NOT NULL DEFAULT 'Reported',
    start_time             TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    end_time               TIMESTAMPTZ,
    problem_at             TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    created_at             TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    created_by             VARCHAR(255),
    updated_at             TIMESTAMPTZ        NOT NULL DEFAULT current_timestamp,
    updated_by             VARCHAR(255),
    deleted_at             TIMESTAMPTZ,
    deleted_by             VARCHAR(255)
);
