package model

type JwtClaimDto struct {
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

type CreateUserDto struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	UserGroupId int    `json:"user_group_id" validate:"required,number"`
}

type UpdateUserDto struct {
	Id          uint64 `json:"id" validate:"required,number"`
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	UserGroupId int    `json:"user_group_id" validate:"required,number"`
}

type DeleteUserDto struct {
	Id uint64 `json:"id" validate:"required,number"`
}

type LoginUserDto struct {
	UserIdentifier string `json:"user_identifier" validate:"required"`
	Password       string `json:"password" validate:"required"`
}
