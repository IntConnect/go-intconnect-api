package model

import "mime/multipart"

type CreateMachineDocumentRequest struct {
	Code         string                `form:"code" validate:"required"`
	Name         string                `form:"name" validate:"required"`
	Description  string                `form:"description" validate:"required"`
	DocumentFile *multipart.FileHeader `form:"document_file" validate:"required,requiredFile,fileExtension=.png .jpg"`
}
