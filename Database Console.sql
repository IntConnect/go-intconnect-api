-- This script is used to create the database for the Go IntConnect system.
DROP DATABASE IF EXISTS go_intconnect_system;
CREATE DATABASE go_intconnect_system;

INSERT INTO nodes (id, type, label, description, help_text, color, icon, component_name, default_config, created_by)
VALUES (1,
        'input',
        'MQTT Source Input',
        'Node to receive data from MQTT broker',
        'This node connects to an MQTT broker and receives data from a specified topic.',
        '#FF5733',
        'tabler-database-heart',
        'MqttInModal',
        '{
            "broker": "tcp://localhost:1883",
            "topic": "sensor/data",
            "client_id": "mqtt_in_node_1"
        }',
        'system');


-- Node 2: JSON Parser

INSERT INTO nodes (id, type, label, description, help_text, color, icon, component_name, default_config, created_by)
VALUES (2,
        'processor',
        'JSON Parser',
        'Node to parse JSON data',
        'This node parses incoming JSON data and prepares it for further processing.',
        '#33FF57',
        'tabler-database-heart',
        'JsonParserModal',
        '{}',
        'system');


-- Node 3: MQTT Out


INSERT INTO nodes (id, type, label, description, help_text, color, icon, component_name, default_config, created_by)
VALUES (3,
        'output',
        'MQTT Sink Output',
        'Node to send data to MQTT broker',
        'This node connects to an MQTT broker and sends processed data to a specified topic.',
        '#3357FF',
        'tabler-database-heart',
        'MqttOutModal',
        '{}',
        'system');

