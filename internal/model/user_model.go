package model

import "go-intconnect-api/internal/trait"

type JwtClaimRequest struct {
	Id          uint64   `json:"id"`
	Email       string   `json:"email"`
	Username    string   `json:"username"`
	Name        string   `json:"name"`
	Permissions []string `json:"-"`
	RoleId      uint64   `json:"role_id" mapstructure:"role_id"`
	RoleName    string   `json:"role_name" mapstructure:"role_name"`
}

type UserResponse struct {
	Id                uint64             `json:"id"`
	Username          string             `json:"username"`
	Name              string             `json:"name"`
	Email             string             `json:"email"`
	AvatarPath        string             `json:"avatar_path"`
	Status            trait.UserStatus   `json:"status"`
	RoleResponse      *RoleResponse      `json:"role,omitempty" mapstructure:"role"`
	AuditableResponse *AuditableResponse `json:"auditable_response"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,unique=users;username"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,min=3,max=100,unique=users;email"`
	Password string `json:"password" validate:"required,min=3,max=100,weakPassword"`
	RoleId   uint64 `json:"role_id" validate:"required,gte=0"`
}

type UpdateUserRequest struct {
	Id       uint64 `json:"-" validate:"required,gte=0"`
	Username string `json:"username" validate:"required,min=3,max=50,unique=users;username;Id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,min=3,max=100,unique=users;email;Id"`
	Password string `json:"password,omitempty" validate:"omitempty,min=3,max=100,weakPassword"`
	RoleId   uint64 `json:"role_id" validate:"required,gte=0"`
}

type LoginUserRequest struct {
	UserIdentifier string `json:"user_identifier" validate:"required,min=3,max=100"`
	Password       string `json:"password" validate:"required"`
}

func (userResponse *UserResponse) GetAuditableResponse() *AuditableResponse {
	return userResponse.AuditableResponse
}

func (userResponse *UserResponse) SetAuditableResponse(auditableResponse *AuditableResponse) {
	userResponse.AuditableResponse = auditableResponse
}
