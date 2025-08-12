package helper

import "go-intconnect-api/internal/model"

func WriteSuccess(message string, data interface{}) model.ResponseContractModel {
	return model.ResponseContractModel{
		Status:  true,
		Message: message,
		Data:    &data,
	}
}
