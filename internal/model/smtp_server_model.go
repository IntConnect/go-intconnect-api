package model

type SmtpServerResponse struct {
	Id                uint64             `json:"id"`
	Host              string             `json:"host"`
	Port              string             `json:"port"`
	Username          string             `json:"username"`
	Password          string             `json:"password"`
	MailAddress       string             `json:"mail_address"`
	MailName          string             `json:"mail_name"`
	IsActive          bool               `json:"is_active"`
	AuditableResponse *AuditableResponse `json:"auditable"`
}

type CreateSmtpServerRequest struct {
	Host        string `json:"host" validate:"required,min=1,max=100"`
	Port        string `json:"port" validate:"required,min=1,max=100"`
	Username    string `json:"username" validate:"required,min=1,max=100"`
	Password    string `json:"password" validate:"required,min=1,max=100"`
	MailAddress string `json:"mail_address" validate:"required,min=1,max=100"`
	MailName    string `json:"mail_name" validate:"required,min=1,max=100"`
	IsActive    bool   `json:"is_active"`
}

type UpdateSmtpServerRequest struct {
	Id          uint64 `json:"-" validate:"required,number,gte=1"`
	Host        string `json:"host" validate:"required,min=1,max=100"`
	Port        string `json:"port" validate:"required,min=1,max=100"`
	Username    string `json:"username" validate:"required,min=1,max=100"`
	Password    string `json:"password" validate:"required,min=1,max=100"`
	MailAddress string `json:"mail_address" validate:"required,min=1,max=100"`
	MailName    string `json:"mail_name" validate:"required,min=1,max=100"`
	IsActive    bool   `json:"is_active"`
}

func (smtpServerResponse *SmtpServerResponse) GetAuditableResponse() *AuditableResponse {
	return smtpServerResponse.AuditableResponse
}

func (smtpServerResponse *SmtpServerResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	smtpServerResponse.AuditableResponse = auditableResponse
}
