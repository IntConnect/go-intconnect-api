package trait

type UserStatus string

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
	StatusBanned   UserStatus = "BANNED"
)
