CREATE TABLE check_sheet_values
(
    id                                         SERIAL PRIMARY KEY,
    check_sheet_id                             BIGINT       NOT NULL REFERENCES check_sheets (id),
    check_sheet_document_template_parameter_id BIGINT       NOT NULL REFERENCES check_sheet_document_templates_parameters (id),
    timestamp                                  VARCHAR(255) NOT NULL,
    value                                      TEXT         NOT NULL
)