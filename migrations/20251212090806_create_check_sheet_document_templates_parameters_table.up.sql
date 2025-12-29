CREATE TABLE check_sheet_document_templates_parameters
(
    id                               SERIAL NOT NULL PRIMARY KEY,
    check_sheet_document_template_id BIGINT NOT NULL REFERENCES check_sheet_document_templates (id),
    parameter_id                     BIGINT NOT NULL REFERENCES parameters (id),
    UNIQUE (check_sheet_document_template_id, parameter_id)
)