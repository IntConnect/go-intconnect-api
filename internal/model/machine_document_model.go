package model

import "mime/multipart"

type CreateMachineDocumentRequest struct {
	Id           *uint64               `form:"id" validate:"omitempty"`
	Name         string                `form:"name" validate:"required"`
	Description  string                `form:"description" validate:"required"`
	DocumentFile *multipart.FileHeader `form:"document_file" validate:"required,requiredFile,fileExtension=.png .jpg"`
}

type MachineDocumentResponse struct {
	Id                uint               `json:"id"`
	MachineId         uint64             `json:"machine_id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	FilePath          string             `json:"file_path"`
	AuditableResponse *AuditableResponse `json:"auditable" mapstructure:"-"`
}

func (machineDocumentResponse *MachineDocumentResponse) GetAuditableResponse() *AuditableResponse {
	return machineDocumentResponse.AuditableResponse
}

func (machineDocumentResponse *MachineDocumentResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	machineDocumentResponse.AuditableResponse = auditableResponse
}
