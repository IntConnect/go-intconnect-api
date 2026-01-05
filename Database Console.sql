-- This script is used to create the database for the Go IntConnect system.
DROP DATABASE IF EXISTS go_intconnect_system;
CREATE DATABASE go_intconnect_system;


INSERT INTO nodes (id, type, label, name, description, help_text, color, icon, component_name, default_config,
                   created_by,
                   deleted_at)
VALUES (1,
        'input',
        'MQTT Source Input',
        'mqtt_input',
        'Node to receive data from MQTT broker',
        'This node connects to an MQTT broker and receives data from a specified topic.',
        '#FF5733',
        'tabler-database-heart',
        'MqttInModal',
        '{
            "action": "single_sub",
            "topic": "sensor/data",
            "qos": "0", "output": "auto-detect-1", "name": "MQTT In"
        }',
        'system',
        NULL);

-- Node 2: JSON Parser

INSERT INTO nodes (id, type, label, name, description, help_text, color, icon, component_name, default_config,
                   created_by,
                   deleted_at)
VALUES (2,
        'processor',
        'JSON Parser',
        'json_parser',
        'Node to parse JSON data',
        'This node parses incoming JSON data and prepares it for further processing.',
        '#33FF57',
        'tabler-database-heart',
        'JsonParserModal',
        '{}',
        'system',
        NULL);


-- Node 3: MQTT Out

INSERT INTO nodes (id, type, label, name, description, help_text, color, icon, component_name, default_config,
                   created_by,
                   deleted_at)
VALUES (3,
        'output',
        'MQTT Sink Output',
        'mqtt_output',
        'Node to send data to MQTT broker',
        'This node connects to an MQTT broker and sends processed data to a specified topic.',
        '#3357FF',
        'tabler-database-heart',
        'MqttOutModal',
        '{}',
        'system',
        NULL);

INSERT INTO nodes (id, type, label, name, description, help_text, color, icon, component_name, default_config,
                   created_by,
                   deleted_at)
VALUES (4,
        'output',
        'Database Sink Output',
        'db_output',
        'Node to send data to database',
        'This node send to an database',
        '#3357FF',
        'tabler-database-heart',
        'DatabaseModal',
        '{}',
        'system',
        NULL);

INSERT INTO mqtt_topics (mqtt_broker_id, name, qos)
VALUES (1, 'sensor/data', 0);
SELECT *
FROM parameters;
INSERT INTO mqtt_brokers(host_name, mqtt_port, ws_port, is_active)
VALUES ('10.175.16.39', '1883', '9001', true);

SELECT *
FROM audit_logs;
SELECT *
FROM machine_documents;
SELECT *
FROM machines;
SELECT *
FROM parameters;
SELECT *
FROM users;
SELECT *
FROM facilities;
SELECT *
FROM machines;
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
FROM parameters WHERE is_watch = true;

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

DELETE
FROM check_sheet_values;
DELETE
FROM check_sheets;
DELETE
FROM parameters;
DELETE
FROM telemetries;
DELETE FROM log_alarms;
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

ALTER TABLE log_alarms
    RENAME COLUMN notes TO note;

SELECT * FROM log_alarms;
DELETE FROM log_alarms;