package auth

import (
	"errors"
	"log"
	"delivery_management_system/internal/db"
	"delivery_management_system/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate checks if the user's credentials are correct
func Authenticate(username, password string) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == db.ErrNoRows {
			return nil, errors.New("user not found")
		}
		log.Println("Error fetching user:", err)
		return nil, err
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	return &user, nil
}
