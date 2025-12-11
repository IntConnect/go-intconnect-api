package model

import "mime/multipart"

type CreateMachineDocumentRequest struct {
	Code         string                `form:"code" validate:"required"`
	Name         string                `form:"name" validate:"required"`
	Description  string                `form:"description" validate:"required"`
	DocumentFile *multipart.FileHeader `form:"document_file" validate:"required,requiredFile,fileExtension=.png .jpg"`
}

type MachineDocumentResponse struct {
	Id                uint               `json:"id"`
	MachineId         uint64             `json:"machine_id"`
	Code              string             `json:"code"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	FilePath          string             `json:"file_path"`
	AuditableResponse *AuditableResponse `json:"auditable" mapstructure:"auditable"`
}

func (machineDocumentResponse *MachineDocumentResponse) GetAuditableResponse() *AuditableResponse {
	return machineDocumentResponse.AuditableResponse
}

func (machineDocumentResponse *MachineDocumentResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	machineDocumentResponse.AuditableResponse = auditableResponse
}
