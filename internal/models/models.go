package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Order struct {
	ID          uuid.UUID              `db:"id" json:"id"`
	UserID      uuid.UUID              `db:"user_id" json:"user_id"`
	Source      string                 `db:"source" json:"source"`
	Destination string                 `db:"destination" json:"destination"`
	Status      string                 `db:"status" json:"status"`
	Metadata    map[string]interface{} `db:"metadata" json:"metadata"`
	CreatedAt   time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at" json:"updated_at"`
}
