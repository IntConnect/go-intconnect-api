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

SELECT * FROM users;
SELECT * FROM nodes;

INSERT INTO users (id, username, email, password) VALUES (1, 'admin', 'admin@gmail.com', '$2a$12$TyWbaQx6LLW8Ik0ZNjCkbeYcKr96XtaTBsJn8uTxmWM/2FRD5AIPG')
UPDATE nodes SET deleted_at = null;
UPDATE users     SET deleted_at = null;
SELECT * FROM pipelines;
INSERT INTO "pipeline_nodes" ("pipeline_id","node_id","type","label","position_x","position_y","created_at","created_by","updated_at","updated_by","deleted_at","deleted_by") VALUES (27,1,'input','MQTT Source Input',191.61653867305168,95.05595775545625,'2025-10-17 23:38:39.868','System','2025-10-17 23:38:39.868','System',NULL,NULL) RETURNING "id"

SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = 'pipeline_nodes';
