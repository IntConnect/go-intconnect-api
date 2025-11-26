CREATE TABLE report_document_templates_parameters
(
    report_document_template_id BIGINT NOT NULL REFERENCES report_document_templates (id),
    parameter_id                BIGINT NOT NULL REFERENCES parameters (id),
    PRIMARY KEY (report_document_template_id, parameter_id)
)