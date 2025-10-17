package model

type JwtClaimRequest struct {
	Email    *string `json:"email" mapstructure:"email"`
	Username *string `json:"username" mapstructure:"username"`
	Role     string  `json:"role" mapstructure:"role"`
}

type UserResponse struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at" mapstructure:"-"`
	UpdatedAt string `json:"updated_at" mapstructure:"-"`
}

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	UserGroupId int    `json:"user_group_id" validate:"required,number"`
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
	UserIdentifier string `json:"user_identifier" validate:"required"`
	Password       string `json:"password" validate:"required"`
}
