CREATE TABLE protocol_configurations
(
    id               SERIAL PRIMARY KEY,
    name             VARCHAR(255) NOT NULL, -- nama koneksi (misalnya "MQTT Broker 1")
    protocol         VARCHAR(255) NOT NULL, -- mqtt, websocket, modbus, dll
    description      TEXT,                  -- opsional, keterangan koneksi
    specific_setting JSONB        NOT NULL, -- konfigurasi disimpan dalam JSON
    created_by       VARCHAR(255),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by       VARCHAR(255),
    deleted_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_by       VARCHAR(255)
)