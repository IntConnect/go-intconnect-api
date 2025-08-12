package model

type JwtClaimDto struct {
	Email    *string `json:"email" mapstructure:"email"`
	Username *string `json:"username" mapstructure:"username"`
	Role     string  `json:"role" mapstructure:"role"`
}
