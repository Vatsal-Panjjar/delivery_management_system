package models

import "time"

type User struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"` // "customer" or "admin"
	CreatedAt    time.Time `db:"created_at"`
}
