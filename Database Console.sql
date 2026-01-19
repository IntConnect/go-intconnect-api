-- This script is used to create the database for the Go IntConnect system.
DROP DATABASE IF EXISTS go_intconnect_system;
CREATE DATABASE go_intconnect_system;

SELECT *
FROM audit_logs;
SELECT *
FROM machine_documents;
SELECT *
FROM machines;
SELECT *
FROM users;
SELECT *
FROM facilities;
SELECT *
FROM mqtt_topics;
SELECT *
FROM system_settings;
SELECT *
FROM parameter_operations;
SELECT *
FROM mqtt_brokers;
SELECT *
FROM telemetries;
SELECT *
FROM parameters;
SELECT *
FROM parameters;
SELECT *
FROM telemetries;
SELECT *
FROM report_document_templates;
SELECT *
FROM roles;
SELECT *
FROM parameters;
SELECT *
FROM machines;
SELECT *
FROM check_sheet_document_templates;
SELECT *
FROM check_sheets;
SELECT *
FROM check_sheet_values;
SELECT *
FROM dashboard_widgets;
DELETE
FROM dashboard_widgets;
SELECT *
FROM modbus_servers;
SELECT *
FROM registers;
SELECT *
FROM dashboard_widgets;
SELECT *
FROM parameters
WHERE is_featured = TRUE;
SELECT *
FROM mqtt_topics;
SELECT *
FROM users;
SELECT *
FROM roles;
SELECT *
FROM roles_permissions;
SELECT *
FROM alarm_logs;
SELECT *
FROM processed_parameter_sequences;
SELECT *
FROM permissions;
DELETE FROM alarm_logs;
DELETE
FROM check_sheet_values;
DELETE
FROM check_sheets;
DELETE
FROM parameters;
DELETE
FROM telemetries;
DELETE
FROM alarm_logs;
DELETE
FROM dashboard_widgets;
DELETE
FROM facilities;
DELETE
FROM machines;

SELECT bucket,
       parameter_id,
       last_value
FROM (SELECT time_bucket_gapfill('5 minutes'::interval, timestamp) AS bucket,
             parameter_id,
             last(value, timestamp)                                AS last_value
      FROM telemetries
      WHERE parameter_id IN (46, 47, 48, 49)
        AND timestamp BETWEEN '2025-12-20 15:04:00' AND '2025-12-20 19:04:00'
      GROUP BY bucket, parameter_id) q
ORDER BY bucket;

ALTER TABLE check_sheet_document_templates
    ADD COLUMN starting_hour TIME NOT NULL default CURRENT_TIMESTAMP;

SELECT *
FROM telemetries
WHERE parameter_id = 1
ORDER BY id DESC;
SELECT *
FROM check_sheet_check_points;
SELECT *
FROM check_sheet_check_point_values;
SELECT * FROM roles_permissions;
SELECT bucket,
       parameter_id,
       last_value
FROM (SELECT time_bucket_gapfill('1 hours'::interval, timestamp, '2026-01-09 08:00:00'::timestamptz,
                                 '2026-01-10 08:00:00'::timestamptz) AS bucket,
             parameter_id,
             last(value, timestamp)                                  AS last_value
      FROM telemetries
      WHERE parameter_id IN (2)
        AND timestamp BETWEEN '2026-01-09 08:00:00' AND '2026-01-10 08:00:00'
      GROUP BY bucket, parameter_id) q
ORDER BY bucket;

ALTER TABLE machine_documents DROP COLUMN code;

DROP TABLE check_sheet_values;

CREATE TABLE users_test
(
    id          SERIAL PRIMARY KEY,
    role_id     BIGINT REFERENCES roles (id) NOT NULL,
    username    VARCHAR(255)                 NOT NULL UNIQUE,
    name        VARCHAR(255)                 NOT NULL UNIQUE,
    email       VARCHAR(255)                 NOT NULL UNIQUE,
    password    VARCHAR(255)                 NOT NULL,
    avatar_path VARCHAR(255)                 NOT NULL,
    created_by  VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP DEFAULT NOW(),
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMP DEFAULT NOW(),
    deleted_by  VARCHAR(255)
)