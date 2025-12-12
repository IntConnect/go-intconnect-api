CREATE TABLE check_sheet_document_templates_parameters
(
    check_sheet_document_template_id BIGINT NOT NULL REFERENCES check_sheet_document_templates (id),
    parameter_id                     BIGINT NOT NULL REFERENCES parameters (id),
    PRIMARY KEY (check_sheet_document_template_id, parameter_id)
)