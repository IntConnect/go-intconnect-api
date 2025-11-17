package model

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required|min:3|max:50"`
	Description string `json:"description" validate:""`
}

type UpdateRoleRequest struct {
	Id          uint64 `json:"id" validate:"required|number"`
	Name        string `json:"name" validate:"required|min:3|max:50"`
	Description string `json:"description" validate:""`
}

type RoleResponse struct {
	Id          uint64 `json:"id" validate:"required|number"`
	Name        string `json:"name" validate:"required|min:3|max:50"`
	Description string `json:"description" validate:""`
}
