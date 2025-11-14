package trait

type UserStatus string

const (
	UserStatusActive   UserStatus = "Active"
	UserStatusInactive UserStatus = "Inactive"
	UserStatusBanned   UserStatus = "Banned"
)
