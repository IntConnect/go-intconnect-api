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
