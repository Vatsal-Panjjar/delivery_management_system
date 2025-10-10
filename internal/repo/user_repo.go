package repo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(user *models.User) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, username, password_hash, role, created_at)
		 VALUES ($1,$2,$3,$4,$5)`,
		user.ID, user.Username, user.PasswordHash, user.Role, time.Now(),
	)
	return err
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Get(&user, "SELECT * FROM users WHERE username=$1", username)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}
