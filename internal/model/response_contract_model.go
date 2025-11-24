package model

type ResponseContractModel struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Error   *interface{} `json:"error,omitempty"`
}

func NewResponseContractModel(statusResp bool, messageResp string, dataResp *interface{}, Error *interface{}) *ResponseContractModel {
	return &ResponseContractModel{
		Status:  statusResp,
		Message: messageResp,
		Data:    dataResp,
		Error:   Error,
	}
}

// BaseResponse adalah struktur dasar untuk semua response
type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SuccessResponse untuk response sukses dengan data
type SuccessResponse[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

type ErrorResponse struct {
	BaseResponse
	Error *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail berisi detail error
type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}
