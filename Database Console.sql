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


DELETE
FROM pipeline_edges;
DELETE
FROM pipeline_nodes;
DELETE
FROM pipelines WHERE id = 5;
DELETE
FROM pipeline_nodes WHERE pipeline_id = 5;