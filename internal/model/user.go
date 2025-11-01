package model

type User struct {
	ID           int64  `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password" json:"-"`
	AuditField
}
