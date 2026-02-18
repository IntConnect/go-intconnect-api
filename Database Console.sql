-- This script is used to create the database for the Go IntConnect system.
DROP DATABASE IF EXISTS go_intconnect_system;
CREATE DATABASE go_intconnect_system;

DELETE FROM dashboard_widgets;;
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

SELECT *
FROM telemetries
WHERE parameter_id = 1
ORDER BY id DESC;
SELECT *
FROM check_sheet_check_points;
SELECT *
FROM check_sheet_check_point_values;
SELECT *
FROM roles_permissions;
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

ALTER TABLE machine_documents
    DROP COLUMN code;

DROP TABLE check_sheet_values;
SELECT *
FROM telemetries
ORDER BY id DESC;

SELECT *
FROM registers;

ALTER TABLE registers
    ADD COLUMN position_x FLOAT NOT NULL DEFAULT 0,
    ADD COLUMN position_y FLOAT NOT NULL DEFAULT 0,
    ADD COLUMN position_z FLOAT NOT NULL DEFAULT 0,
    ADD COLUMN rotation_x FLOAT NOT NULL DEFAULT 0,
    ADD COLUMN rotation_y FLOAT NOT NULL DEFAULT 0,
    ADD COLUMN rotation_z FLOAT NOT NULL DEFAULT 0;


ALTER TABLE check_sheet_document_templates
    ADD COLUMN
        interval_type VARCHAR(255) NOT NULL DEFAULT 'Hour';


SELECT *
FROM users;

SELECT * FROM mqtt_brokers;
DELETE FROM mqtt_brokers WHERE id =11;
DELETE FROM mqtt_topics WHERE mqtt_broker_id != 1 OR mqtt_broker_id != 2;
DELETE FROM parameters WHERE id IN (SELECT id from mqtt_topics WHERE mqtt_broker_id != 1 OR mqtt_broker_id != 2) ;

UPDATE roles SET deleted_at = NULL;

DROP TABLE IF EXISTS "public"."telemetries";
CREATE TABLE "public"."telemetries" (
  "id" SERIAL,
  "parameter_id" int8 NOT NULL,
  "value" float8,
  "timestamp" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP
);

SELECT * FROM users;
SELECT * FROM check_sheet_document_templates;
SELECT * FROM mqtt_brokers;
SELECT * FROM permissions;
SELECT * FROM dashboard_widgets;
DELETE FROM dashboard_widgets WHERE id = 32 OR id = 33;