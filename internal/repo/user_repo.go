package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/Vatsal-Panjiar/delivery_management_system/internal/models"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.DB.Get(&u, "SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Create(u *models.User) error {
	_, err := r.DB.Exec("INSERT INTO users (id, username, password_hash, role, created_at) VALUES ($1,$2,$3,$4,$5)", u.ID, u.Username, u.PasswordHash, u.Role, u.CreatedAt)
	return err
}
