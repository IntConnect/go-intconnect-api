CREATE TYPE machine_parameter_label_type_enum AS ENUM ('hover', 'always', 'panel');

CREATE TABLE machines_parameters
(
    id           SERIAL PRIMARY KEY,
    parameter_id BIGINT                            NOT NULL REFERENCES parameters (id),
    mesh_name    VARCHAR(255)                      NOT NULL,
    position_x   FLOAT                             NOT NULL,
    position_y   FLOAT                             NOT NULL,
    position_z   FLOAT                             NOT NULL,
    label_type   machine_parameter_label_type_enum NOT NULL
)