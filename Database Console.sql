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

SELECT *
FROM users;
SELECT *
FROM nodes;

INSERT INTO users (id, username, email, password)
VALUES (1, 'admin', 'admin@gmail.com', '$2a$12$TyWbaQx6LLW8Ik0ZNjCkbeYcKr96XtaTBsJn8uTxmWM/2FRD5AIPG')

SELECT *
FROM pipelines;
SELECT *
FROM pipeline_nodes;
SELECT *
FROM pipeline_edges;
SELECT *
FROM database_connections;
SELECT *
FROM nodes;
SELECT *
FROM facilities;
SELECT *
FROM roles;
SELECT *
FROM mqtt_brokers;

DELETE
FROM pipeline_edges;
DELETE
FROM pipeline_nodes;
DELETE
FROM pipelines;
DELETE
FROM database_connections;



CREATE DATABASE sensors;
SELECT current_database() AS database_name,
       c.table_name,
       c.column_name,
       c.data_type,
       c.is_nullable,
       c.column_default
FROM information_schema.columns c
         JOIN information_schema.tables t
              ON c.table_name = t.table_name
WHERE t.table_schema = 'public'
ORDER BY c.table_name, c.ordinal_position;

SELECT *
FROM mqtt_topics;

SELECT *
FROM telemetries;

SELECT *
FROM machines;


SELECT *
FROM permissions;

INSERT INTO mqtt_topics (mqtt_broker_id, name, qos)
VALUES (1, 'sensor/data', 0);

CREATE EXTENSION IF NOT EXISTS timescaledb;
\dx