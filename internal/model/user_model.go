package model

type JwtClaimRequest struct {
	Email    *string `json:"email" mapstructure:"email"`
	Username *string `json:"username" mapstructure:"username"`
	Role     string  `json:"role" mapstructure:"role"`
}

type UserResponse struct {
	Id                uint64             `json:"id"`
	Username          string             `json:"username"`
	Name              string             `json:"name"`
	Email             string             `json:"email"`
	AuditableResponse *AuditableResponse `json:"auditable_response"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,unique=users;username;id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,min=3,max=100"`
	Password string `json:"password" validate:"required,min=3,max=100"`
	RoleId   uint64 `json:"role_id" validate:"required,gte=0"`
}

type UpdateUserRequest struct {
	Id          uint64 `json:"id" validate:"required,number"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	UserGroupId int    `json:"user_group_id" validate:"required,number"`
}

type DeleteUserRequest struct {
	Id uint64 `json:"id" validate:"required,number"`
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
