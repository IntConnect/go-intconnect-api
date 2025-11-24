package model

type DeleteResourceGeneralRequest struct {
	Id uint64 `json:"id" validate:"required|number"`
}
