CREATE TABLE breakdowns
(
    id                     SERIAL PRIMARY KEY,
    machine_id             BIGINT       NOT NULL REFERENCES machines (id),
    reported_by            BIGINT       NOT NULL REFERENCES users (id),
    verified_by            BIGINT REFERENCES users (id),

    problem_identification TEXT         NOT NULL DEFAULT '',
    people_issue           TEXT         NOT NULL DEFAULT '',
    inspection_issue       TEXT         NOT NULL DEFAULT '',
    condition_issue        TEXT         NOT NULL DEFAULT '',
    action_taken           TEXT         NOT NULL DEFAULT '',
    parts_issue            TEXT         NOT NULL DEFAULT '',
    analysis_notes         TEXT         NOT NULL DEFAULT '',

    corrective_action      TEXT,
    preventive_action      TEXT,
    status                 VARCHAR(255) NOT NULL DEFAULT 'Reported',
    start_time             TIMESTAMP  NOT NULL DEFAULT current_timestamp,
    end_time               TIMESTAMP,
    problem_at             TIMESTAMP  NOT NULL DEFAULT current_timestamp,
    created_at             TIMESTAMP  NOT NULL DEFAULT current_timestamp,
    created_by             VARCHAR(255),
    updated_at             TIMESTAMP  NOT NULL DEFAULT current_timestamp,
    updated_by             VARCHAR(255),
    deleted_at             TIMESTAMP,
    deleted_by             VARCHAR(255)
);
