CREATE TYPE machine_parameter_label_type_enum AS ENUM ('hover', 'always', 'panel');

CREATE TABLE machines_parameters
(
    id           SERIAL PRIMARY KEY,
    parameter_id BIGINT                            NOT NULL REFERENCES parameters (id),
    mesh_name    VARCHAR(255)                      NOT NULL,
    pos_x        FLOAT                             NOT NULL,
    pos_y        FLOAT                             NOT NULL,
    pos_z        FLOAT                             NOT NULL,
    label_type   machine_parameter_label_type_enum NOT NULL
)