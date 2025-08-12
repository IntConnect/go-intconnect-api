package entity

type User struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Auditable
}
